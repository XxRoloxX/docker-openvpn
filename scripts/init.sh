OVPN_DATA="openvpn-data"

docker build . -t rolo-openvpn

docker volume create --name $OVPN_DATA

IP_ADDRESS="172.104.141.218"

echo $IP_ADDRESS

docker run -v "$OVPN_DATA:/etc/openvpn" --rm rolo-openvpn ovpn_genconfig -u "udp://$IP_ADDRESS"

docker run -v "$OVPN_DATA:/etc/openvpn" --rm -it rolo-openvpn ovpn_initpki


