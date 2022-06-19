// SPDX-License-Identifier: MIT
import { bold, italic, parse, underline } from "./deps.ts";
import { QueryResponse, ServerOptions } from "./lib/utils.ts";
import { doQuery } from "./lib/query.ts";
import { parseResponse } from "./lib/response.ts";
import { parseArgs } from "./args.ts";
import { parseNAPTR, parsePTR } from "./lib/reverse.ts";

async function main() {
  // Parse args
  const args = parse(Deno.args, {
    alias: {
      h: "help",
      p: "port",
      s: "short",
      x: "ptr",
      V: "version",
    },
    boolean: ["help", "ptr", "short", "version"],
    string: ["port"],
    default: {
      "port": "53",
    },
  });
  if (args.help) args.version = true;

  if (args.version) {
    console.log(
      `${
        bold("awl")
      } version 0.1.0 (running with deno ${Deno.version.deno}, TypeScript ${Deno.version.typescript}, on V8 ${Deno.version.v8})
Written by (YOUR NAME GOES HERE)`,
    );
  }

  if (args.help) {
    console.log(
      ` ${bold("Usage:")} awl name ${italic("type")} ${italic("@server")}
       ${bold("<name>")}    domain name
       ${bold("<type>")}    defaults to A
       ${bold("<@server>")} defaults to your local resolver

       Order ${bold("DOES NOT")} matter\n
      `,
      `${underline("Options")}:
        -p <port> use <port> for query, defaults to 53
        -s Equivalent to dig +short, only return addresses
        -x do a reverse lookup
        -h print this helpful guide
        -V get the version
      `,
    );
  }

  if (args.version) Deno.exit(0);

  const parsedArgs = parseArgs(args);

  let domain = parsedArgs.name || ".";

  let query = parsedArgs.type || "A";
  if (domain === ".") query = "NS";

  if (query === "PTR") {
    // The "server" is an IP address, it needs to become a canonical domain
    domain = parsePTR(domain);
  } else if (query === "NAPTR") {
    // The "server" is a phone number, it needs to become a canonical domain
    domain = parseNAPTR(domain);
  }

  const server = parsedArgs.server || { server: "", port: 53 };

  const response: QueryResponse = await doQuery(domain, query, server);

  if (!args.short) {
    console.log(
      `;; ->>HEADER<<- opcode: QUERY, rcode: ${response.response}
;; QUESTION SECTION:
;; ${domain}    IN    ${query}
  
;; ANSWER SECTION:`,
    );
  }
  if (response.response === "NOERROR") {
    const res = parseResponse(response, domain, query, args.short);
    res.forEach((answer) => {
      console.log(answer);
    });
  }
  if (!args.short) {
    console.log(`
;; ADDITONAL SECTION:

;; Query time: ${response.time} msec
;; SERVER: ${displayServer(server)}
;; WHEN: ${new Date().toLocaleString()}
        `);
  }
}

/**
 * A handler for displaying the server
 * @param {ServerOptions} server - The DNS server used
 * @returns {string} The string used
 */
function displayServer(server: ServerOptions): string {
  let val = "";
  if (server.server) {
    val += server.server;
  }
  if (server.port != 53) {
    val += `#${server.port}`;
  }
  if (!val) val = "System";
  return val;
}

if (import.meta.main) {
  await main();
}
