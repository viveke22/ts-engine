// Type definitions for ts-engine built-ins

declare module "http" {
    export interface IncomingMessage {
        url: string;
        method: string;
    }

    export interface ServerResponse {
        writeHead(statusCode: number, headers?: { [key: string]: string }): void;
        end(data?: string | any): void;
    }

    export interface Server {
        listen(port: number, callback?: () => void): void;
    }

    export function createServer(requestListener: (req: IncomingMessage, res: ServerResponse) => void): Server;
}

declare function require(moduleName: string): any;

// Global fetch support
declare function fetch(url: string): {
    status: number;
    ok: boolean;
    statusText: string;
};

// Console support
interface Console {
    log(...args: any[]): void;
}
declare var console: Console;
