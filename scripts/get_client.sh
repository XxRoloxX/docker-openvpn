CLIENTNAME=$1
OVPN_DATA="openvpn-data"

docker run -v "$OVPN_DATA:/etc/openvpn" --rm -it rolo-openvpn easyrsa build-client-full "$CLIENTNAME" nopass


docker run -v "$OVPN_DATA:/etc/openvpn" --rm kylemanna/openvpn ovpn_getclient "$CLIENTNAME" > "$CLIENTNAME.ovpn"
