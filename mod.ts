// SPDX-License-Identifier: MIT

// Exports for lawl, the library for awl
export type { QueryResponse, ServerOptions } from "./lib/utils.ts";
export { isRecordType } from "./lib/utils.ts";
export { doQuery } from "./lib/query.ts";
export { parseResponse } from "./lib/response.ts";
export { parseIPv6, parseNAPTR, parsePTR } from "./lib/reverse.ts";
