package com.alipay.application.service.account.cloud.sensecore;

import java.util.List;

import com.alipay.application.service.account.cloud.Credential;
/*
 *@title SensecoreCredential
 *@description
 *@author Center-Sun
 *@version 1.0
 *@create 2026/1/17 17:39
 */
public class SensecoreCredential implements Credential{

    private final String ak;
    private final String sk;

    public String getAk() {
        return ak;
    }

    public String getSk() {
        return sk;
    }

    public SensecoreCredential(String ak, String sk) {
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
