package com.alipay.application.service.account.cloud.ucloud;

import com.alipay.application.service.account.cloud.Credential;
import java.util.List;
/*
 *@title UCloudCredential
 *@description
 *@author Center-Sun
 *@version 1.0
 *@create 2025/8/30 17:39
 */
public class UCloudCredential implements Credential{
    private final String ak;
    private final String sk;

    public String getAk() {
        return ak;
    }

    public String getSk() {
        return sk;
    }

    public UCloudCredential(String ak, String sk) {
        this.ak = ak;
        this.sk = sk;
    }

    @Override
    public boolean verification() {
        try {
            regions();
        } catch (Exception e) {
            throw new RuntimeException("Cloud account verification failed:" + e.getMessage());
        }
        return true;
    }

    @Override
    public List<Region> regions() {
        return List.of();
    }
}
