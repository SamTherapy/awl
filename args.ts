
import { Args } from "./deps.ts";
import { isRecordType, ServerOptions } from "./lib/utils.ts";
/**
 * A handler for parsing the arguments passed in
 * @param {ServerOptions} server - The DNS server to query
 * @param {Deno.RecordType} type - The type of DNS request, see Deno.RecordType for more info
 * @param {string} name - Server to look up
 */
export type arguments = {
  server?: ServerOptions;
  type?: Deno.RecordType;
  name?: string;
};

/**
 * @param {Args} args - The arguments, directly passed in
 * @returns {arguments} The arguments, parsed
 */
export function parseArgs(args: Args): arguments {
  const parsed: arguments = {} as arguments;
  args._.forEach((arg) => {
    arg = arg.toString();

    // if it starts with an @, it's a server
    if (arg.includes("@")) {
      parsed.server = {
        server: arg.split("@").pop() as string,
        port: args.port,
      };
      return;
    }
    // if there is a dot, it's a name
    if (arg.includes(".")) {
      parsed.name = arg;
      return;
    }

    if (isRecordType(arg)) {
      parsed.type = arg.toUpperCase() as Deno.RecordType;
      return;
    }

    // if all else fails, assume it's a name
    parsed.name = arg;
  });

  // Add a . to the end of the name if it's not there
  if (parsed.name?.charAt(parsed.name.length - 1) !== ".") {
    parsed.name = parsed.name?.concat(".");
  }

  return parsed;
}
