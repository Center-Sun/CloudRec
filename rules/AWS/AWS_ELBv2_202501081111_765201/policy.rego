package elb_sg_rule_open_to_all_public_1900001_286
import rego.v1

default risk := false

risk if {
    sg_rule_open_to_all
    pub_lb
    has_listener
}

## 基础信息
load_balancer_name := input.ELB.LoadBalancerName
dns_name := input.ELB.DNSName

## elb 网络类型
net_scheme := input.ELB.Scheme

## 安全组入向规则
sg_rules contains ip_permission if {
    some ip_permission in input.SecurityGroups[_].SecurityGroup.IpPermissions
}

sg_rule_open_to_all if {
    some ip_permission in sg_rules
    sg_rule_open_to_all_public(ip_permission)
    sg_rule_allows_listener(ip_permission)
}
sg_rule_open_to_all if {
    ## 不存在安全组
    not input.SecurityGroups
}

sg_rule_open_to_all_public(ip_permission) if {
    some ip_ranges in ip_permission.IpRanges
    cidr_ip := ip_ranges.CidrIp
    cidr_ip == "0.0.0.0/0"
}

sg_rule_allows_listener(ip_permission) if {
    ip_permission.IpProtocol == "-1"
}

sg_rule_allows_listener(ip_permission) if {
    some port in listener_ports
    from_port := ip_permission.FromPort
    to_port := ip_permission.ToPort
    is_number(from_port)
    is_number(to_port)
    from_port <= port
    port <= to_port
}

has_listener if {
    count(listener_ports) > 0
}

listener_ports contains port if {
    some listener in input.Listeners
    port := listener.Port
    is_number(port)
}

listener_ports contains port if {
    some listener in input.Listeners
    port := listener.Listener.Port
    is_number(port)
}

## 公网lb
pub_lb if {
    net_scheme in ["internet-facing"]
}
