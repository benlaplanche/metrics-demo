# Metrics Demo Release

A BOSH release to demo deploying a single node Redis server with the metron agent located on the same VM, by taking the Metron Agent from tjhe CF Release itself. 

This demo release will query the `redis info` command to obtain the `uptime_in_seconds` value and emit this as a metric. 

## Latest Release

**Version:** [http://metrics-demo-blobs.s3.amazonaws.com/metrics-demo-2.tgz](http://metrics-demo-blobs.s3.amazonaws.com/metrics-demo-2.tgz)

Tested with **CF Release 213**

## Deploy to BOSH-Lite

```
$ bosh upload release http://metrics-demo-blobs.s3.amazonaws.com/metrics-demo-2.tgz
$ bosh deployment manifests/bosh-lite.yml
$ bosh deploy -n
```

## Consume the emitted metrics

```
$ bosh deployment manifests/bosh-lite.yml
$ bosh ssh
$ sudo -i
```

Now you can execute the binary to consume the emitted metrics
```
/var/vcap/packages/consumer/bin/metrics-consumer --config=/var/vcap/jobs/consumer/config.json
```

You will get output similar to this
```
**Started consuming the firehose**
origin:"metrics-demo/z1/0" eventType:ValueMetric timestamp:1439218008244245766 deployment:"cf-warden" job:"redis" index:"0" ip:"10.244.0.118" valueMetric:<name:"numCPUS" value:8 unit:"count" >
origin:"metrics-demo/z1/0" eventType:ValueMetric timestamp:1439218008244417884 deployment:"cf-warden" job:"redis" index:"0" ip:"10.244.0.118" valueMetric:<name:"memoryStats.numBytesAllocatedStack" value:278528 unit:"count" >
origin:"metrics-demo/z1/0" eventType:ValueMetric timestamp:1439218008244465044 deployment:"cf-warden" job:"redis" index:"0" ip:"10.244.0.118" valueMetric:<name:"memoryStats.lastGCPauseTimeNS" value:131322 unit:"count" >
origin:"metrics-demo/z1/0" eventType:ValueMetric timestamp:1439218008244454786 deployment:"cf-warden" job:"redis" index:"0" ip:"10.244.0.118" valueMetric:<name:"memoryStats.numFrees" value:392 unit:"count" >
origin:"metrics-demo/z1/0" eventType:ValueMetric timestamp:1439218008244593437 deployment:"cf-warden" job:"redis" index:"0" ip:"10.244.0.118" valueMetric:<name:"uptime_in_seconds" value:292 unit:"" >
```

The first lines are free courtesy of Metron.

The metric we have chosen to emit is this one, the last entry
```
origin:"metrics-demo/z1/0" eventType:ValueMetric timestamp:1439218008244593437 deployment:"cf-warden" job:"redis" index:"0" ip:"10.244.0.118" valueMetric:<name:"uptime_in_seconds" value:292 unit:"" >
```

Metrics are emitted on a 15second interval. 

## Config changes

Default values are set in the `spec` files of the Jobs.

If you need to change these values you can override them in your manifest.

The most commonly changed ones would look like this in your manifest
```
jobs:
- name: redis
  instances: 1
  resource_pool: metrics-demo
  networks:
  - name: default
  templates:
  - name: redis
    release: metrics-demo
  - name: emitter
    release: metrics-demo
  - name: consumer
    release: metrics-demo
  - name: metron_agent
    release: cf
  properties:
    metron_agent:
      zone: z1
      debug: false
    consumer:
      uaa_url: https://uaa.xxx
      doppler_addr: wss://doppler.xxx:4443
      username: xxxx
      password: xxxx
```