import { CreateServer } from "./components/create-server";
import { Toaster } from "./components/ui/sonner";
// import { Welcome } from "./components/welcome";

export function App() {
	return (
		<div className="min-h-screen bg-background text-foreground flex items-center justify-center">
			<CreateServer />
			<Toaster />
		</div>
	);
}

export default App;
