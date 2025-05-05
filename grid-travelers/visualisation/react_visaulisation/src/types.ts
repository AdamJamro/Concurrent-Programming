export interface GridParams {
    cols: number;
    rows: number;
}

export interface TravelerUpdate {
    [id: string]: {
        x: number;
        y: number;
    };
}

export const isGridParams = (message: WebSocketMessage): message is GridParams => {
    return (message as GridParams).cols !== undefined && (message as GridParams).rows !== undefined;
}

export type WebSocketMessage = | GridParams | TravelerUpdate;