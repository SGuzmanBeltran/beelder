import { Card, CardContent } from "./ui/card";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "./ui/select";

import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { Slider } from "./ui/slider";
import { Switch } from "./ui/switch";

export function CreateServer() {
	return (
		<div className="flex flex-col items-center justify-center min-h-150 space-y-8 w-2/3">
			<div className="flex w-full justify-start">
				<h1 className=" text-2xl font-bold">Create your server</h1>
			</div>

			<div className="flex w-full space-x-8">
				<div className="flex flex-col w-4/5 space-y-8">
					<Card className="w-full">
						<CardContent className="space-y-4">
							<h3 className="text-lg">Configure your installation</h3>
							<div>
								<h4>What should we install on you server?</h4>
								<Select>
									<SelectTrigger className="w-full">
										<SelectValue placeholder="Select server type" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="vanilla">Vanilla</SelectItem>
										<SelectItem value="paper">Paper</SelectItem>
										<SelectItem value="forge">Forge</SelectItem>
									</SelectContent>
								</Select>

								<Select>
									<SelectTrigger className="w-full">
										<SelectValue placeholder="Select server version" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="1.21.11">1.21.11</SelectItem>
										<SelectItem value="1.21.10">1.21.10</SelectItem>
										<SelectItem value="1.21.9">1.21.9</SelectItem>
									</SelectContent>
								</Select>
							</div>
							<div>
								<h4>
									How many players do you expect to be on the server at once?
								</h4>
								<Slider defaultValue={[0]} max={100} step={1} />
							</div>

							<div>
								<h4>This is the recommended configuration for your server</h4>
								<Card className="w-full bg-stone-800">
									{/* Show the recommended configuration based on the previous inputs and pricing*/}
								</Card>
							</div>
						</CardContent>
					</Card>

					<Card className="w-full">
						<CardContent>
							<h3 className="text-lg">Select your location</h3>
							<Select>
								<SelectTrigger className="w-full">
									<SelectValue placeholder="Select a region" />
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="us-east-1">
										US East (N. Virginia)
									</SelectItem>
									<SelectItem value="us-west-2">US West (Oregon)</SelectItem>
									<SelectItem value="eu-west-1">EU West (Ireland)</SelectItem>
									<SelectItem value="ap-southeast-1">
										Asia Pacific (Singapore)
									</SelectItem>
									<SelectItem value="sa-east-1">
										South America (São Paulo)
									</SelectItem>
								</SelectContent>
							</Select>
						</CardContent>
					</Card>

					<Card className="w-full pb-8">
						<CardContent className="space-y-4">
							<h3 className="text-lg">Add your initial configuration</h3>
							<div>
								<h4>What's your server name?</h4>
								<Input />
							</div>

							<div className="flex justify-between">
								<div className="w-1/2 pr-3">
									<h4>What difficulty?</h4>
									<Select>
										<SelectTrigger className="w-full">
											<SelectValue placeholder="Theme" />
										</SelectTrigger>
										<SelectContent>
											<SelectItem value="peaceful">Peaceful</SelectItem>
											<SelectItem value="easy">Easy</SelectItem>
											<SelectItem value="normal">Normal</SelectItem>
											<SelectItem value="hard">Hard</SelectItem>
											<SelectItem value="hardcore">Hardcore</SelectItem>
										</SelectContent>
									</Select>
								</div>
								<div className="w-1/2 flex-col justify-center items-center">
									<h4 className="pb-2">Do we allow only premium player?</h4>
									<Switch />
								</div>
							</div>
						</CardContent>
					</Card>
				</div>
				<div className="flex flex-col w-2/5 space-y-8">
					<Card className="max-w-2xl w-full">
						<CardContent className="space-y-6">
							<h3 className="text-2xl font-bold">Order Summary</h3>

							<div className="space-y-2">
								<div className="flex justify-between items-start">
									<div>
										<h4 className="font-semibold">Server Hosting Package</h4>
										<p className="text-sm text-muted-foreground">
											Paper 1.21.11 - 6GB RAM
										</p>
									</div>
									<span className="font-semibold">$22.49/mo</span>
								</div>

								<div className="text-sm space-y-1">
									<p>
										<span className="text-muted-foreground">
											Server Location:
										</span>{" "}
										Miami, Florida
									</p>
								</div>
							</div>

							<div className="border-t pt-4 space-y-2">
								<h4 className="font-semibold">Totals</h4>
								<div className="flex justify-between">
									<span>Monthly</span>
									<span className="font-semibold">$22.49</span>
								</div>
							</div>

							<div className="border-t pt-4">
								<p className="text-center text-3xl font-bold mb-4">$22.49</p>
								<p className="text-center text-sm text-muted-foreground mb-4">
									Total Due Today
								</p>
							</div>

							<Button className="w-full bg-primary py-6 rounded-md text-sm font-semibold cursor-pointer">
								Create Server
							</Button>

							<p className="text-center text-sm text-muted-foreground cursor-pointer">
								Have a promo code?
							</p>
						</CardContent>
					</Card>
				</div>
			</div>
		</div>
	);
}
