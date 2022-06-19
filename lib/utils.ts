// SPDX-License-Identifier: MIT

/**
 * A DNS response
 */
export type QueryResponse = {
  dnsResponse:
    | string[]
    | Deno.CAARecord[]
    | Deno.MXRecord[]
    | Deno.NAPTRRecord[]
    | Deno.SOARecord[]
    | Deno.SRVRecord[]
    | string[][];
  response: string;
  time: number;
};
/**
 * Options for which DNS server to query
 */
export type ServerOptions = {
  server: string;
  port?: number;
};

export function isRecordType(type: string): type is Deno.RecordType {
  return type.toUpperCase() === "A" || type.toUpperCase() === "AAAA" || type.toUpperCase() === "CNAME" || type.toUpperCase() === "MX" ||
    type.toUpperCase() === "NS" || type.toUpperCase() === "PTR" || type.toUpperCase() === "SOA" || type.toUpperCase() === "TXT" ||
    type.toUpperCase() === "NAPTR" || type.toUpperCase() === "SRV" || type.toUpperCase() === "CAA";
}

/**
 * Test if the DNS query is an MX record
 * @param {QueryResponse["dnsResponse"]} record - DNS response
 * @returns {boolean} - true if the record is an MX record
 */
export function isMX(
  record: QueryResponse["dnsResponse"],
): record is Deno.MXRecord[] {
  return (record as unknown as Deno.MXRecord[])[0].exchange !== undefined;
}
/**
 * Test if the DNS query is a CAA record
 * @param {QueryResponse["dnsResponse"]} record - DNS response
 * @returns {boolean} - true if the record is a CAA record
 */
export function isCAA(
  record: QueryResponse["dnsResponse"],
): record is Deno.CAARecord[] {
  return (record as unknown as Deno.CAARecord[])[0].critical !== undefined;
}

/**
 * Test if the DNS query is an NAPTR record
 * @param {QueryResponse["dnsResponse"]} record - DNS response
 * @returns {boolean} - true if the record is an NAPTR record
 */
export function isNAPTR(
  record: QueryResponse["dnsResponse"],
): record is Deno.NAPTRRecord[] {
  return (record as unknown as Deno.NAPTRRecord[])[0].regexp !== undefined;
}

/**
 * Test if the DNS query is an SOA record
 * @param {QueryResponse["dnsResponse"]} record - DNS response
 * @returns {boolean} - true if the record is an SOA record
 */
export function isSOA(
  record: QueryResponse["dnsResponse"],
): record is Deno.SOARecord[] {
  return (record as unknown as Deno.SOARecord[])[0].rname !== undefined;
}

/**
 * Test if the DNS query is an SRV record
 * @param {QueryResponse["dnsResponse"]} record - DNS response
 * @returns {boolean} - true if the record is an SRV record
 */
export function isSRV(
  record: QueryResponse["dnsResponse"],
): record is Deno.SRVRecord[] {
  return (record as unknown as Deno.SRVRecord[])[0].port !== undefined;
}

export function isTXT(
  record: QueryResponse["dnsResponse"],
): record is string[][] {
  return Array.isArray(record as unknown as Array<string>);
}
