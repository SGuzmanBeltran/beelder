import "./index.css";

import App from "./App.tsx";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

// Set title based on environment
const stage = import.meta.env.VITE_STAGE;
if (stage && stage !== "production") {
	document.title = `Beelder - ${stage}`;
}

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<App />
	</StrictMode>
);
