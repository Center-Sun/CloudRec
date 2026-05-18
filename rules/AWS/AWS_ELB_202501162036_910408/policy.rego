package ob_aws_elbv2_unsafe_protocol_5300003_285

import rego.v1

default risk := false

risk if {
    count(unsafe_protocol) > 0
}

elb_id := input.ELB.LoadBalancerArn

unsafe_protocol[p] if {
    some p in input.Listeners
    p.Protocol == "HTTP"
}

unsafe_protocol[p] if {
    some p in input.Listeners
    p.Listener.Protocol == "HTTP"
}
