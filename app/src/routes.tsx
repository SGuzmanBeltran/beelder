import { CreateServer } from "./pages/create-server";
import { RootLayout } from "./layouts/RootLayout";
import { Welcome } from "./pages/welcome";
import { ServerDetails } from "./pages/server-details";
import {
  createBrowserRouter,
} from "react-router-dom";

export const router = createBrowserRouter([
	{
		path: "/",
		element: <RootLayout />,
		children: [
			{ index: true, element: <Welcome /> },
			{ index: false, path: "create-server", element: <CreateServer /> },
			{ index: false, path: "server/:id", element: <ServerDetails /> },
			{ index: false, path: "servers", element: <div>Servers List</div> },
			{ path: "*", element: <div>Not Found</div> },
		],
	},
]);