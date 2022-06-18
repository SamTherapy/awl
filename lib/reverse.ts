// SPDX-License-Identifier: MIT

export function parsePTR(ip: string) {
  if (ip.includes(".")) {
    // It's an IPv4 address
    const ptr = ip.split(".");
    let pop: string | undefined = "not undefined";
    let domain = "";
    do {
      pop = ptr.pop();
      if (pop) {
        domain += `${pop}.`;
      }
    } while (pop !== undefined);
    domain += "in-addr.arpa";
    return domain;
  } else if (ip.includes(":")) {
    const parsedIP = parseIPv6(ip);
    const ptr = parsedIP.split(":");
    // It's an IPv6 address
    let pop: string[] | undefined = ["e"];
    let domain = "";
    do {
      pop = ptr.pop()?.split("").reverse();
      if (pop) {
        for (const part of pop) {
          domain += `${part}.`;
        }
      }
    } while (pop !== undefined);
    domain += "ip6.arpa";
    return domain;
  } else {
    // It's not an address
    return "";
  }
}

export function parseIPv6(addr: string) {
  addr = addr.replace(/^:|:$/g, "");

  const ipv6 = addr.split(":");

  for (let i = 0; i < ipv6.length; i++) {
    let hex: string | string[] = ipv6[i];
    if (hex != "") {
      // normalize leading zeros
      // TODO: make this not deprecated
      ipv6[i] = ("0000" + hex).substr(-4);
    } else {
      // normalize grouped zeros ::
      hex = [];
      for (let j = ipv6.length; j <= 8; j++) {
        hex.push("0000");
      }
      ipv6[i] = hex.join(":");
    }
  }

  return ipv6.join(":");
}

export function parseNAPTR(phNum: string) {
  phNum = phNum.toString();
  phNum = phNum.replace("+", "").replaceAll(" ", "").replaceAll("-", "");
  const rev = phNum.split("").reverse();
  let ptr = "";
  rev.forEach((n) => {
    ptr += `${n}.`;
  });
  ptr += "e164.arpa";
  return ptr;
}
