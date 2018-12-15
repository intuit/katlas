
export function validateIPaddress(ipaddress) {
    const ipformat = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
    if (ipaddress.match(ipformat)) {
        //alert(`Valid Ip address: ${ipaddress}`);
        return true;
    }
    //alert(`Invalid Ip address: ${ipaddress}`);
    return false;
}
