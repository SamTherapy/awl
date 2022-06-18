// SPDX-License-Identifier: MIT
import { parseIPv6, parseNAPTR, parsePTR } from "../lib/reverse.ts";
import { assertEquals } from "./testDeps.ts";

Deno.test("IPv6 Parse, localhost", () => {
  assertEquals(parseIPv6("::1"), "0000:0000:0000:0000:0000:0000:0000:0001");
});

Deno.test("IPv6 Parse, :: in middle of address", () => {
  assertEquals(
    parseIPv6("2001:4860:4860::8844"),
    "2001:4860:4860:0000:0000:0000:0000:8844",
  );
});

Deno.test("IPv4 PTR, localhost", () => {
  assertEquals(parsePTR("127.0.0.1"), "1.0.0.127.in-addr.arpa");
});

Deno.test("IPv6 PTR, actual IP", () => {
  assertEquals(
    parsePTR("2606:4700:4700::1111"),
    "1.1.1.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.7.4.0.0.7.4.6.0.6.2.ip6.arpa",
  );
});

Deno.test("PTR, Fallback lel", () => {
  assertEquals(parsePTR("1367218g3a1"), "");
});

Deno.test("NAPTR, US number", () => {
  assertEquals(
    parseNAPTR("+1-800-555-1234"),
    "4.3.2.1.5.5.5.0.0.8.1.e164.arpa",
  );
});

Deno.test("NAPTR, non-US number", () => {
  assertEquals(
    parseNAPTR("44 186 533 2244"),
    "4.4.2.2.3.3.5.6.8.1.4.4.e164.arpa",
  );
});
