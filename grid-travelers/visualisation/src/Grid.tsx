import { Stage, Line, Layer, Rect, Circle, Text } from 'react-konva';
import useWebSocket from 'react-use-websocket';
import React from "react";
import { GridParams, WebSocketMessage } from './types';

const GRID_STROKE_COLOR = '#ddd';
// const GRID_BACKGROUND_COLOR = '#aa6644';
const CHECKER_COLOR_1 = '#ff00ff'; // Light color
const CHECKER_COLOR_2 = '#a000a0'; // Dark color

const isGridParams = (message: WebSocketMessage): message is GridParams => {
    return (message as GridParams).cols !== undefined && (message as GridParams).rows !== undefined;
}

export default function Grid({ rows, cols }: { rows: number; cols: number }) {
    const [gridParams, setGridParams] = React.useState<{ rows: number; cols: number }>({ rows, cols });
    const [gridParamsReceived, setGridParamsReceived] = React.useState(false);
    const [travelers, setTravelers] = React.useState<Record<string, { x: number; y: number }>>({});
    const [cellSize, setCellSize] = React.useState(30);

    const { lastMessage } = useWebSocket('ws://localhost:8080/ws', {
        shouldReconnect: () => true,
    });

    React.useEffect(() => {
        if (lastMessage?.data) {
            const data = JSON.parse(lastMessage.data);
            console.log('Received message:', data);
            if (isGridParams(data)) {
                setGridParams(data);
                setCellSize(
                    Math.min(window.innerWidth / gridParams.cols, window.innerHeight / gridParams.rows) * 0.9
                );
                setGridParamsReceived(true);
            } else if(gridParamsReceived) {
                setTravelers(prev => ({ ...prev, ...data }));
            } else {
                console.error('Received unexpected message:', data);
            }
        }
    }, [gridParamsReceived, lastMessage]);
    
    // React.useEffect(() => {
    //     const width = window.innerWidth;
    //     const height = window.innerHeight;
    //     const cellSize = Math.min(width / gridParams.cols, height / gridParams.rows) * 0.9;
    //     setCellSize(cellSize);
    // }, [gridParams]);

    // if (!gridParamsReceived) {
    //     return <div>Loading...</div>;
    // }

    return (
        <Stage className={"grid-stage"} width={gridParams.cols * cellSize} height={gridParams.rows * cellSize}>
            <Layer>
                {/* Checker-like Background */}
                {Array.from({ length: gridParams.rows }).map((_, row) =>
                    Array.from({ length: gridParams.cols }).map((_, col) => (
                        <Rect
                            key={`${row}-${col}`}
                            x={col * cellSize}
                            y={row * cellSize}
                            width={cellSize}
                            height={cellSize}
                            fill={(row + col) % 2 === 0 ? CHECKER_COLOR_1 : CHECKER_COLOR_2}
                            stroke={GRID_STROKE_COLOR}
                            strokeWidth={1}
                        />
                    ))
                )}

                {/* Vertical Lines */}
                {Array.from({ length: gridParams.cols + 1 }).map((_, i) => (
                    <Line
                        key={`vline-${i}`}
                        points={[i * cellSize, 0, i * cellSize, rows * cellSize]}
                        stroke={GRID_STROKE_COLOR}
                        strokeWidth={1}
                    />
                ))}

                {/* Horizontal Lines */}
                {Array.from({ length: gridParams.rows + 1 }).map((_, i) => (
                    <Line
                        key={`hline-${i}`}
                        points={[0, i * cellSize, gridParams.cols * cellSize, i * cellSize]}
                        stroke={GRID_STROKE_COLOR}
                        strokeWidth={1}
                    />
                ))}

                {/* Travelers */}
                {Object.entries(travelers).map(([id, traveler]) => (
                    <>
                    <Circle
                        key={id}
                        x={(traveler.x + 0.5) * cellSize}
                        y={(traveler.y + 0.5) * cellSize}
                        radius={cellSize / 3}
                        fill={'#332233'}
                    />

                    {/* Label for the traveler */}
                    <Text
                        key={`label-${id}`}
                        x={(traveler.x + 0.5) * cellSize - cellSize / 6} // Center the label horizontally
                        y={(traveler.y + 0.5) * cellSize - cellSize / 5}
                        text={id} // Use the traveler's ID as the label
                        fontSize={12}
                        fill={'#FFF'}
                        align="center"
                    />
                    </>
                ))}


            </Layer>
        </Stage>
    );
}