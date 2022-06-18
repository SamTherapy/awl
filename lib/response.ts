// SPDX-License-Identifier: MIT
import {
  isCAA,
  isMX,
  isNAPTR,
  isSOA,
  isSRV,
  isTXT,
  QueryResponse,
} from "./utils.ts";

/**
 * @param res A DNS {@link QueryResponse}
 * @param domain The domain (or IP, if doing it in reverse) queried
 * @param query The DNS query that was queried
 * @returns {string} The DNS response, put in canonical form
 */
export function parseResponse(
  res: QueryResponse,
  domain: string,
  query: string,
  short: boolean,
): string[] {
  const answer: string[] = [];
  switch (query) {
    case "A":
    case "AAAA":
    case "CNAME":
    case "NS":
    case "PTR":
      res.dnsResponse.forEach((ip) => {
        let dnsQuery = "";
        if (!short) dnsQuery += `${domain}       IN    ${query}    `;
        dnsQuery += `${ip}`;
        answer.push(dnsQuery);
      });
      break;
    case "MX":
      if (isMX(res.dnsResponse)) {
        res.dnsResponse.forEach((record) => {
          let dnsQuery = "";
          if (!short) {
            dnsQuery += `${domain}       IN    ${query}    `;
          }
          dnsQuery += `${record.preference} ${record.exchange}`;
          answer.push(dnsQuery);
        });
      }
      break;
    case "CAA":
      if (isCAA(res.dnsResponse)) {
        res.dnsResponse.forEach((record) => {
          let dnsQuery = "";
          if (!short) {
            dnsQuery += `${domain}       IN    ${query}    `;
          }
          dnsQuery += `${
            record.critical ? "1" : "0"
          } ${record.tag} "${record.value}"`;
          answer.push(dnsQuery);
        });
      }
      break;
    case "NAPTR":
      if (isNAPTR(res.dnsResponse)) {
        res.dnsResponse.forEach((record) => {
          let dnsQuery = "";
          if (!short) dnsQuery += `${domain}       IN    ${query}    `;
          dnsQuery +=
            `${record.order} ${record.preference} "${record.flags}" "${record.services}" ${record.regexp} ${record.replacement}`;
          answer.push(dnsQuery);
        });
      }
      break;
    case "SOA":
      if (isSOA(res.dnsResponse)) {
        res.dnsResponse.forEach((record) => {
          let dnsQuery = "";
          if (!short) dnsQuery += `${domain}       IN    ${query}    `;
          dnsQuery +=
            `${record.mname} ${record.rname} ${record.serial} ${record.refresh} ${record.retry} ${record.expire} ${record.minimum}`;
          answer.push(dnsQuery);
        });
      }
      break;
    case "SRV":
      if (isSRV(res.dnsResponse)) {
        res.dnsResponse.forEach((record) => {
          let dnsQuery = "";
          if (!short) dnsQuery += `${domain}       IN    ${query}    `;
          dnsQuery +=
            `${record.priority} ${record.weight} ${record.port} ${record.target}`;
          answer.push(dnsQuery);
        });
      }
      break;
    case "TXT":
      if (isTXT(res.dnsResponse)) {
        res.dnsResponse.forEach((record) => {
          let dnsQuery = "";
          let txt = "";
          record.forEach((value) => {
            txt += `"${value}"`;
          });
          if (!short) {
            dnsQuery += `${domain}       IN    ${query}    `;
          }
          dnsQuery += `${txt}`;
          answer.push(dnsQuery);
        });
      }
      break;
    default:
      throw new Error("Not yet implemented");
  }
  return answer;
}
