// SPDX-License-Identifier: MIT
import { assertEquals } from "./testDeps.ts";
import { doQuery } from "../lib/query.ts";

Deno.test("Get localhost", async () => {
  const res = await doQuery("localhost", "A");
  assertEquals(res.dnsResponse, ["127.0.0.1"]);
  assertEquals(res.response, "NOERROR");
});

Deno.test("Get localhost, external NS", async () => {
  const res = await doQuery("localhost", "AAAA", { server: "1.1.1.1" });
  assertEquals(res.dnsResponse, ["::1"]);
  assertEquals(res.response, "NOERROR");
});

Deno.test("PTR localhost", async () => {
  const res = await doQuery("1.0.0.127.in-addr.arpa.", "PTR");
  assertEquals(res.dnsResponse, ["localhost."]);
});

// This test will fail if this random Bri ish phone number goes down
// It's also unreliable, so it's disabled
// Deno.test("NAPTR, Remote",async () => {
//   const res = await doQuery("4.4.2.2.3.3.5.6.8.1.4.4.e164.arpa.", "NAPTR");
//   assertStrictEquals(res.dnsResponse, [
//     {
//       order: 100,
//       preference: 10,
//       flags: "u",
//       services: "E2U+sip",
//       regexp: "!^\\+441865332(.*)$!sip:\\1@nominet.org.uk!",
//       replacement: "."
//     },
//     {
//       order: 100,
//       preference: 20,
//       flags: "u",
//       services: "E2U+pstn:tel",
//       regexp: "!^(.*)$!tel:\\1!",
//       replacement: "."
//     }
//   ])
// })

Deno.test("Get invalid IP, regular NS", async () => {
  const res = await doQuery("l", "A");
  assertEquals(res.dnsResponse, undefined);
  assertEquals(res.response, "NXDOMAIN");
});

// This isn't supposed to SERVFAIL
// It also takes forever

Deno.test("Get invalid IP, external NS", async () => {
  const res = await doQuery("b", "AAAA", { server: "1.1.1.1" });
  assertEquals(res.dnsResponse, undefined);
  assertEquals(res.response, "SERVFAIL");
});
