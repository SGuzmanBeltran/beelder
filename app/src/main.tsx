import "./index.css";

import { RouterProvider } from "react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { router } from "./routes.tsx";

// Set title based on environment
const stage = import.meta.env.VITE_STAGE;
if (stage && stage !== "production") {
	document.title = `Beelder - ${stage}`;
}

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<RouterProvider router={router} />,
	</StrictMode>,
);
