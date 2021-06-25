## nri-ss
Publish metrics from [iproute2's](https://wiki.linuxfoundation.org/networking/iproute2) `ss` ([Socket Statistics](https://git.kernel.org/pub/scm/network/iproute2/iproute2.git/tree/misc/ss.c)) command which is the modern replacement for `netstat`.

`iproute2` and `ss` _should_ be preinstalled on most current Linux distributions.

### Installation
1. Place the `nri-ss` executable in `/var/db/newrelic-infra/custom-integrations/`
2. Place `nri-ss-definition.yml` in `/var/db/newrelic-infra/custom-integrations/`
3. Place `nri-ss-config.yml` in `/etc/newrelic-infra/integrations.d/`

### Configuration
- `resolve` : boolean (true | false). If true attemp to resolve host ip addresses to names
  - _IMPORTANT_ : if `resolve` is `true` then filters _must_ use host names. If `resolve` is `false` then filters must use host ip addresses.
- `filter`  : a _properly_ formatted `ss` filter string. No validation is performed, test from the command line directly with `ss`
  - No filter (empty filter) retrieves all available metrics
  - Examples
    - `'( dst 1.2.3.4 )'` show only metrics where the destination host is ip address `1.2.3.4` (`resolve: false`)
    - `'( dst www.google.com )'` show only metrics where the distination host is `www.google.com` (`resolve: true`)
    - `'( dst www.google.com:https )'` show only metrics where the distination host is `www.google.com` and the port is 443 (`resolve: true`)
    - `'( dst 1.2.3.4 or dst 5.6.7.8 )'` show only metrics where the destination host is ip address `1.2.3.4` _or_ the destination host is ip address `5.6.7.8` (`resolve: false`)
    - Filters can be quite complex, perform an Internet search on `ss filter examples` to learn more
- `ss_args` : command line arguments to `ss`, the default is `-iot`

### Trouble shooting
- Try your filter from the command line `ss -iot <your_filter_here>`

### Building
```bash
# Create your root dir
# cd to your root dir
export GOPATH=`pwd`
mkdir -p src/github.com/newrelic/
cd `src/github.com/newrelic/
git clone <this_repo>
cd nri-ss
govendor fetch +
make
```

### Metrics
`"eventType": "SocketStatisticsSample"`
```
              ts     show string "ts" if the timestamp option is set

              sack   show string "sack" if the sack option is set

              ecn    show string "ecn" if the explicit congestion notification option is set

              ecnseen
                     show string "ecnseen" if the saw ecn flag is found in received packets

              fastopen
                     show string "fastopen" if the fastopen option is set

              cong_alg
                     the congestion algorithm name, the default congestion algorithm is "cubic"

              wscale:<snd_wscale>:<rcv_wscale>
                     if window scale option is used, this field shows the send scale factory and receive scale factory

              rto:<icsk_rto>
                     tcp re-transmission timeout value, the unit is millisecond

              backoff:<icsk_backoff>
                     used for exponential backoff re-transmission, the actual re-transmission timeout value is icsk_rto << icsk_backoff

              rtt:<rtt>/<rttvar>
                     rtt is the average round trip time, rttvar is the mean deviation of rtt, their units are millisecond

              ato:<ato>
                     ack timeout, unit is millisecond, used for delay ack mode

              mss:<mss>
                     max segment size

              cwnd:<cwnd>
                     congestion window size

              pmtu:<pmtu>
                     path MTU value

              ssthresh:<ssthresh>
                     tcp congestion window slow start threshold

              bytes_acked:<bytes_acked>
                     bytes acked

              bytes_received:<bytes_received>
                     bytes received

              segs_out:<segs_out>
                     segments sent out

              segs_in:<segs_in>
                     segments received

              send <send_bps>bps
                     egress bps

              lastsnd:<lastsnd>
                     how long time since the last packet sent, the unit is millisecond

              lastrcv:<lastrcv>
                     how long time since the last packet received, the unit is millisecond

              lastack:<lastack>
                     how long time since the last ack received, the unit is millisecond

              pacing_rate <pacing_rate>bps/<max_pacing_rate>bps
                     the pacing rate and max pacing rate

              rcv_space:<rcv_space>
                     a helper variable for TCP internal auto tuning socket receive buffer
```
