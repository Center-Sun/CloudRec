// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blb

import (
	"context"
	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/cloudrec/baidu/collector"
	"github.com/core-sdk/constant"
	"github.com/core-sdk/log"
	"github.com/core-sdk/schema"
	"go.uber.org/zap"
)

type AppBLBDetail struct {
	AppBLB                   appblb.AppBLBModel
	ListenerList             []appblb.AppAllListenerModel
	SecurityGroups           []appblb.BlbSecurityGroupModel
	EnterpriseSecurityGroups []appblb.BlbEnterpriseSecurityGroupModel
	AppServerGroupList       []AppServerGroupDetail
	AppIpGroupList           []AppIpGroupDetail
}

// AppServerGroupDetail enriches an app server group with its backend real
// servers (DescribeBlbRs), which the group description endpoint does not return.
// The embedded AppServerGroup is anonymous so its fields stay flattened and
// existing JSON paths (e.g. $.AppServerGroupList[*].id) keep resolving.
type AppServerGroupDetail struct {
	appblb.AppServerGroup
	BackendServerList []appblb.AppBackendServer `json:"backendServerList"`
}

// AppIpGroupDetail enriches an app IP group with its backend members
// (DescribeAppIpGroupMember). An APP BLB uses either server groups or IP groups
// for its backends; collecting both makes the real backends visible regardless.
type AppIpGroupDetail struct {
	appblb.AppIpGroup
	MemberList []appblb.AppIpGroupMember `json:"memberList"`
}

func GetAppBLBResource() schema.Resource {
	return schema.Resource{
		ResourceType:      collector.APPBLB,
		ResourceTypeName:  "APP BLB",
		ResourceGroupType: constant.NET,
		Desc:              `https://cloud.baidu.com/doc/BLB/s/Lkcznyjer`,
		Regions: []string{
			"blb.bj.baidubce.com",
			"blb.gz.baidubce.com",
			"blb.su.baidubce.com",
			"blb.hkg.baidubce.com",
			"blb.fwh.baidubce.com",
			"blb.bd.baidubce.com",
			"blb.fsh.baidubce.com",
			"blb.sin.baidubce.com",
		},
		ResourceDetailFunc: func(ctx context.Context, service schema.ServiceInterface, res chan<- any) error {
			client := service.(*collector.Services).APPBLBClient

			args := &appblb.DescribeLoadBalancersArgs{}
			for {
				response, err := client.DescribeLoadBalancers(args)
				if err != nil {
					log.CtxLogger(ctx).Warn("DescribeLoadBalancers error", zap.Error(err))
					return err
				}
				for _, i := range response.BlbList {
					d := AppBLBDetail{
						AppBLB:                   i,
						ListenerList:             describeAppAllListeners(ctx, client, i.BlbId),
						SecurityGroups:           describeAppBLBSecurityGroups(ctx, client, i.BlbId),
						EnterpriseSecurityGroups: describeAppBLBEnterpriseSecurityGroups(ctx, client, i.BlbId),
						AppServerGroupList:       describeAppServerGroup(ctx, client, i.BlbId),
						AppIpGroupList:           describeAppIpGroup(ctx, client, i.BlbId),
					}
					res <- d
				}
				if response.NextMarker == "" {
					break
				}
				args.Marker = response.NextMarker
			}

			return nil
		},
		RowField: schema.RowField{
			ResourceId:   "$.AppBLB.blbId",
			ResourceName: "$.AppBLB.name",
			Address:      "$.AppBLB.publicIp",
		},
		Dimension: schema.Regional,
	}
}

func describeAppAllListeners(ctx context.Context, client *appblb.Client, blbId string) (listenerList []appblb.AppAllListenerModel) {
	args := &appblb.DescribeAppListenerArgs{
		Marker:  "",
		MaxKeys: 50,
	}

	for {
		response, err := client.DescribeAppAllListeners(blbId, args)
		if err != nil {
			log.CtxLogger(ctx).Warn("DescribeAppAllListeners error", zap.Error(err))
			return
		}
		listenerList = append(listenerList, response.ListenerList...)
		if response.NextMarker == "" {
			break
		}
		args.Marker = response.NextMarker
	}

	return listenerList
}

func describeAppBLBSecurityGroups(ctx context.Context, client *appblb.Client, blbId string) []appblb.BlbSecurityGroupModel {
	resp, err := client.DescribeSecurityGroups(blbId)
	if err != nil {
		log.CtxLogger(ctx).Warn("describeAppBLBSecurityGroups error", zap.Error(err))
		return nil
	}

	return resp.BlbSecurityGroups
}

func describeAppBLBEnterpriseSecurityGroups(ctx context.Context, client *appblb.Client, blbId string) []appblb.BlbEnterpriseSecurityGroupModel {
	resp, err := client.DescribeEnterpriseSecurityGroups(blbId)
	if err != nil {
		log.CtxLogger(ctx).Warn("describeAppBLBEnterpriseSecurityGroups error", zap.Error(err))
		return nil
	}

	return resp.BlbEnterpriseSecurityGroups
}

func describeAppServerGroup(ctx context.Context, client *appblb.Client, blbId string) (appServerGroupList []AppServerGroupDetail) {
	args := &appblb.DescribeAppServerGroupArgs{
		Marker:  "",
		MaxKeys: 50,
	}
	for {
		response, err := client.DescribeAppServerGroup(blbId, args)
		if err != nil {
			log.CtxLogger(ctx).Warn("describeAppServerGroup error", zap.Error(err))
			return
		}
		for _, sg := range response.AppServerGroupList {
			appServerGroupList = append(appServerGroupList, AppServerGroupDetail{
				AppServerGroup:    sg,
				BackendServerList: describeBlbRs(ctx, client, blbId, sg.Id),
			})
		}
		if response.NextMarker == "" {
			break
		}
		args.Marker = response.NextMarker
	}

	return appServerGroupList
}

// describeBlbRs lists the backend real servers bound to one app server group.
func describeBlbRs(ctx context.Context, client *appblb.Client, blbId string, sgId string) (backendServerList []appblb.AppBackendServer) {
	args := &appblb.DescribeBlbRsArgs{
		Marker:  "",
		MaxKeys: 50,
		SgId:    sgId,
	}
	for {
		response, err := client.DescribeBlbRs(blbId, args)
		if err != nil {
			log.CtxLogger(ctx).Warn("describeBlbRs error", zap.Error(err))
			return
		}
		backendServerList = append(backendServerList, response.BackendServerList...)
		if response.NextMarker == "" {
			break
		}
		args.Marker = response.NextMarker
	}

	return backendServerList
}

func describeAppIpGroup(ctx context.Context, client *appblb.Client, blbId string) (appIpGroupList []AppIpGroupDetail) {
	args := &appblb.DescribeAppIpGroupArgs{
		Marker:  "",
		MaxKeys: 50,
	}
	for {
		response, err := client.DescribeAppIpGroup(blbId, args)
		if err != nil {
			log.CtxLogger(ctx).Warn("describeAppIpGroup error", zap.Error(err))
			return
		}
		for _, ig := range response.AppIpGroupList {
			appIpGroupList = append(appIpGroupList, AppIpGroupDetail{
				AppIpGroup: ig,
				MemberList: describeAppIpGroupMembers(ctx, client, blbId, ig.Id),
			})
		}
		if response.NextMarker == "" {
			break
		}
		args.Marker = response.NextMarker
	}

	return appIpGroupList
}

// describeAppIpGroupMembers lists the backend members bound to one app IP group.
func describeAppIpGroupMembers(ctx context.Context, client *appblb.Client, blbId string, ipGroupId string) (memberList []appblb.AppIpGroupMember) {
	args := &appblb.DescribeAppIpGroupMemberArgs{
		Marker:    "",
		MaxKeys:   50,
		IpGroupId: ipGroupId,
	}
	for {
		response, err := client.DescribeAppIpGroupMember(blbId, args)
		if err != nil {
			log.CtxLogger(ctx).Warn("describeAppIpGroupMember error", zap.Error(err))
			return
		}
		memberList = append(memberList, response.MemberList...)
		if response.NextMarker == "" {
			break
		}
		args.Marker = response.NextMarker
	}

	return memberList
}
