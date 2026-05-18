package elb_sg_rule_open_to_all_public_1900001_286

import rego.v1

test_internet_facing_open_sg_without_listener_is_not_risk if {
	not risk with input as {
		"ELB": {
			"Scheme": "internet-facing",
		},
		"SecurityGroups": [
			{
				"SecurityGroup": {
					"IpPermissions": [
						{
							"FromPort": 80,
							"ToPort":   80,
							"IpRanges": [
								{
									"CidrIp": "0.0.0.0/0",
								},
							],
						},
					],
				},
			},
		],
		"Listeners": null,
	}
}

test_internet_facing_open_sg_with_listener_is_risk if {
	risk with input as {
		"ELB": {
			"Scheme": "internet-facing",
		},
		"SecurityGroups": [
			{
				"SecurityGroup": {
					"IpPermissions": [
						{
							"FromPort": 80,
							"ToPort":   80,
							"IpRanges": [
								{
									"CidrIp": "0.0.0.0/0",
								},
							],
						},
					],
				},
			},
		],
		"Listeners": [
			{
				"ListenerArn": "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/lb/123/456",
				"Port":        80,
			},
		],
	}
}

test_internet_facing_sg_open_to_unlistened_port_is_not_risk if {
	not risk with input as {
		"ELB": {
			"Scheme": "internet-facing",
		},
		"SecurityGroups": [
			{
				"SecurityGroup": {
					"IpPermissions": [
						{
							"FromPort": 22,
							"ToPort":   22,
							"IpRanges": [
								{
									"CidrIp": "0.0.0.0/0",
								},
							],
						},
					],
				},
			},
		],
		"Listeners": [
			{
				"ListenerArn": "arn:aws:elasticloadbalancing:us-east-1:123456789012:listener/app/lb/123/456",
				"Port":        443,
			},
		],
	}
}
