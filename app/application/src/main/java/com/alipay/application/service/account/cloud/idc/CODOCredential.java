package com.alipay.application.service.account.cloud.idc;

import java.util.List;

import com.alipay.application.service.account.cloud.Credential;
/*
 *@title CODOCredential
 *@description
 *@author Center-Sun
 *@version 1.0
 *@create 2025/11/8 13:39
 */
public class CODOCredential implements Credential{

    private final String ak;
    private final String sk;

    public String getAk() {
        return ak;
    }

    public String getSk() {
        return sk;
    }

    public CODOCredential(String ak, String sk) {
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