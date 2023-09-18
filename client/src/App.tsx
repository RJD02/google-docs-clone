import { useEffect, useState } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";

function App() {
    const [count, setCount] = useState(0);
    useEffect(() => {
        // const wss = new WebSocket("ws://localhost:8000/ws/2");
        // wss.onmessage = (msg) => console.log(msg);
        // wss.send("Hello World");
        const ws = new WebSocket("ws://localhost:8000/ws/2");
        ws.onopen = () => {
            console.log("WebSocket connection established");
            ws.send("Hello");
        };
        ws.onmessage = (event) => console.log(event.data);

        ws.onerror = (err) => console.error(err);

        ws.onclose = () => console.log("Websocket connection closed");

        return () => {
            ws.close();
        };
    }, []);

    return (
        <>
            <div>
                <a href="https://vitejs.dev" target="_blank">
                    <img src={viteLogo} className="logo" alt="Vite logo" />
                </a>
                <a href="https://react.dev" target="_blank">
                    <img src={reactLogo} className="logo react" alt="React logo" />
                </a>
            </div>
            <h1>Vite + React</h1>
            <div className="card">
                <button onClick={() => setCount((count) => count + 1)}>
                    count is {count}
                </button>
                <p>
                    Edit <code>src/App.tsx</code> and save to test HMR
                </p>
            </div>
            <p className="read-the-docs">
                Click on the Vite and React logos to learn more
            </p>
        </>
    );
}

export default App;
