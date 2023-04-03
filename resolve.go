package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func resolve(hostname string) {

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, "8.8.8.8:53")
		},
	}
	ips, _ := r.LookupHost(context.Background(), hostname)

	for _, ip := range ips {
		fmt.Println(ip)
	}

}
