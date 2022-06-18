// SPDX-License-Identifier: MIT
import { assertEquals, assertThrows } from "./testDeps.ts";
import { parseResponse } from "../lib/response.ts";
import { QueryResponse } from "../lib/utils.ts";

const mockResponse: QueryResponse = {
  dnsResponse: [],
  response: "NOERROR",
  time: 0,
};
let domain = "localhost.";

Deno.test("A query", () => {
  mockResponse.dnsResponse = ["127.0.0.1"];
  assertEquals(parseResponse(mockResponse, domain, "A", false), [
    "localhost.       IN    A    127.0.0.1",
  ]);
});

Deno.test("AAAA query, short", () => {
  mockResponse.dnsResponse = ["::1"];
  assertEquals(parseResponse(mockResponse, domain, "AAAA", true), [
    "::1",
  ]);
});

Deno.test("MX query", () => {
  mockResponse.dnsResponse = [{
    exchange: "mail.localhost",
    preference: 10,
  }];
  assertEquals(parseResponse(mockResponse, domain, "MX", false), [
    "localhost.       IN    MX    10 mail.localhost",
  ]);
});

Deno.test("CAA query", () => {
  mockResponse.dnsResponse = [{
    critical: false,
    tag: "issue",
    value: "pki.goog",
  }];
  assertEquals(parseResponse(mockResponse, domain, "CAA", false), [
    'localhost.       IN    CAA    0 issue "pki.goog"',
  ]);
});

Deno.test("NAPTR query", () => {
  domain = "4.3.2.1.5.5.5.0.0.8.1.e164.arpa.";
  mockResponse.dnsResponse = [{
    flags: "u",
    order: 100,
    preference: 10,
    services: "E2U+sip",
    regexp: "!^.*$!sip:customer-service@example.com!",
    replacement: ".",
  }, {
    flags: "u",
    order: 102,
    preference: 10,
    services: "E2U+email",
    regexp: "!^.*$!mailto:information@example.com!",
    replacement: ".",
  }];
  assertEquals(parseResponse(mockResponse, domain, "NAPTR", false), [
    `4.3.2.1.5.5.5.0.0.8.1.e164.arpa.       IN    NAPTR    100 10 "u" "E2U+sip" !^.*$!sip:customer-service@example.com! .`,
    `4.3.2.1.5.5.5.0.0.8.1.e164.arpa.       IN    NAPTR    102 10 "u" "E2U+email" !^.*$!mailto:information@example.com! .`,
  ]);
});

Deno.test("SOA query", () => {
  domain = "cloudflare.com.";
  mockResponse.dnsResponse = [{
    mname: "ns3.cloudflare.com.",
    rname: "dns.cloudflare.com.",
    serial: 2280958559,
    refresh: 10000,
    retry: 2400,
    expire: 604800,
    minimum: 300,
  }];
  assertEquals(parseResponse(mockResponse, domain, "SOA", false), [
    "cloudflare.com.       IN    SOA    ns3.cloudflare.com. dns.cloudflare.com. 2280958559 10000 2400 604800 300",
  ]);
});

Deno.test("SRV query", () => {
  domain = "localhost";
  mockResponse.dnsResponse = [{
    port: 22,
    priority: 0,
    target: "localhost",
    weight: 10,
  }];
  assertEquals(parseResponse(mockResponse, domain, "SRV", false), [
    "localhost       IN    SRV    0 10 22 localhost",
  ]);
});

Deno.test("TXT query", () => {
  mockResponse.dnsResponse = [["a"]];
  assertEquals(parseResponse(mockResponse, domain, "TXT", false), [
    'localhost       IN    TXT    "a"',
  ]);
});

Deno.test("Invalid query", () => {
  mockResponse.dnsResponse = [["a"]];
  assertThrows((): void => {
    parseResponse(mockResponse, domain, "E", true);
  });
});
