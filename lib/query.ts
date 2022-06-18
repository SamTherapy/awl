// SPDX-License-Identifier: MIT
import { QueryResponse, ServerOptions } from "./utils.ts";

/**
 * @param domain Domain to query
 * @param query
 * @param server {@link utils.ts/ServerOptions}
 * @returns
 */
export async function doQuery(
  domain: string,
  query: Deno.RecordType,
  server?: ServerOptions,
) {
  const response: QueryResponse = {} as QueryResponse;
  if (!server?.server) {
    const t0 = performance.now();
    await Deno.resolveDns(domain, query)
      // If there's no error
      .then((value) => {
        const t1 = performance.now();
        response.time = t1 - t0;
        response.response = "NOERROR";
        response.dnsResponse = value;
      })
      // If there is an error
      .catch((e: Error) => {
        const t1 = performance.now();
        response.time = t1 - t0;
        switch (e.name) {
          case "NotFound":
            response.response = "NXDOMAIN";
            break;
          default:
            response.response = "SERVFAIL";
        }
      });
  } else {
    const t0 = performance.now();
    await Deno.resolveDns(domain, query, {
      nameServer: { "ipAddr": server.server, "port": server.port },
    })
      // If there's no error
      .then((value) => {
        const t1 = performance.now();
        response.time = t1 - t0;
        response.response = "NOERROR";
        response.dnsResponse = value;
      })
      // If there is an error
      .catch((e: Error) => {
        const t1 = performance.now();
        response.time = t1 - t0;
        switch (e.name) {
          case "NotFound":
            response.response = "NXDOMAIN";
            break;
          default:
            response.response = "SERVFAIL";
        }
      });
  }
  return response;
}
