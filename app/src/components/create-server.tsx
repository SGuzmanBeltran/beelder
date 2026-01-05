import { ArrowLeft, ArrowRight } from "@hugeicons/core-free-icons";
import { Card, CardContent } from "./ui/card";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "./ui/select";
import { type SubmitHandler, useForm, type FieldErrors } from "react-hook-form";

import { HugeiconsIcon } from "@hugeicons/react";
import { Input } from "./ui/input";
import { PricingCard } from "./pricing-card";
import { Slider } from "./ui/slider";
import { Switch } from "./ui/switch";
import { useEffect, useState } from "react";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

const pricingPlans = [
	{
		ram: "1GB",
		price: "$0.00",
		badge: { text: "Free", color: "stone" as const },
	},
	{
		ram: "2GB",
		price: "$11.99",
	},
	{
		ram: "4GB",
		price: "$15.99",
	},
	{
		ram: "6GB",
		price: "$17.99",
	},
	{
		ram: "8GB",
		price: "$23.99",
	},
	{
		ram: "12GB",
		price: "$29.99",
	},
];

const serverConfigSchema = z.object({
	serverType: z
		.string("Server type is required")
		.min(1, "Server type is required"),
	serverVersion: z
		.string("Server version is required")
		.min(1, "Server version is required"),
	playerCount: z
		.number()
		.min(1, "Player count must be at least 1")
		.max(100, "Player count cannot exceed 100"),
	region: z.string("Region is required").min(1, "Region is required"),
	ramPlan: z.string().min(1, "RAM plan is required"),
	serverName: z
		.string("Server name is required")
		.min(1, "Server name is required"),
	difficulty: z
		.string("Difficulty is required")
		.min(1, "Difficulty is required"),
	premiumOnly: z.boolean(),
});

type ServerConfig = z.infer<typeof serverConfigSchema>;

export function CreateServer() {
	const [currentPlanIndex, setCurrentPlanIndex] = useState(0);
	const [isLoadingRecommendation, setIsLoadingRecommendation] = useState(false);
	const [recommendedPlan, setRecommendedPlan] = useState<number | null>(null);

	const handlePrevious = () => {
		setCurrentPlanIndex((prev) => {
			const newIndex = prev > 0 ? prev - 1 : pricingPlans.length - 1;
			setValue("ramPlan", pricingPlans[newIndex].ram);
			return newIndex;
		});
	};

	const handleNext = () => {
		setCurrentPlanIndex((prev) => {
			const newIndex = prev < pricingPlans.length - 1 ? prev + 1 : 0;
			setValue("ramPlan", pricingPlans[newIndex].ram);
			return newIndex;
		});
	};

	// Get recommended plan from backend
	const fetchRecommendedPlan = async (
		serverType: string,
		playerCount: number,
		region: string
	) => {
		if (!serverType || !playerCount || !region) return;

		setIsLoadingRecommendation(true);
		try {
			const response = await fetch(
				`http://localhost:3000/api/v1/server/recommended-plans?serverType=${serverType}&playersCount=${playerCount}&region=${region}`,
				{
					method: "GET",
				}
			);
			const data = await response.json();

			// Find the index of the recommended plan
			const recommendedRam = data.data.recommendation;
			const planIndex = pricingPlans.findIndex(
				(plan) => plan.ram === recommendedRam
			);
			console.log(planIndex);

			if (planIndex !== -1 && planIndex !== 0) {
				// Don't auto-select free plan
				setCurrentPlanIndex(planIndex);
				setRecommendedPlan(planIndex);
				setValue("ramPlan", pricingPlans[planIndex].ram);
			}
		} catch (error) {
			console.error("Error fetching recommendation:", error);
		} finally {
			setIsLoadingRecommendation(false);
		}
	};

	const {
		handleSubmit,
		watch,
		setValue,
		formState: { errors },
	} = useForm<ServerConfig>({
		resolver: zodResolver(serverConfigSchema),
		defaultValues: {
			playerCount: 0,
			premiumOnly: true,
			ramPlan: pricingPlans[0].ram,
		},
	});

	// Watch for changes in serverType, playerCount, and region
	useEffect(() => {
		const subscription = watch((value, { name }) => {
			if (
				(name === "serverType" ||
					name === "playerCount" ||
					name === "region") &&
				value.serverType &&
				value.playerCount &&
				value.region
			) {
				// Debounce API call
				const timeoutId = setTimeout(() => {
					fetchRecommendedPlan(
						value.serverType!,
						value.playerCount!,
						value.region!
					);
				}, 500);
				return () => clearTimeout(timeoutId);
			}
		});
		return () => subscription.unsubscribe();
	}, [watch]);

	const onSubmit: SubmitHandler<ServerConfig> = (data) =>
		console.log("OnSubmit " + JSON.stringify(data));

	const onError = (errors: FieldErrors<ServerConfig>) => {
		console.log("Form Errors: " + JSON.stringify(errors));
	};

	// Determine badge for current plan
	const getCurrentPlanBadge = () => {
		if (currentPlanIndex === 0) {
			return { text: "Free", color: "stone" as const };
		} else if (currentPlanIndex === recommendedPlan) {
			return { text: "Recommended", color: "yellow" as const };
		} else if (recommendedPlan && currentPlanIndex < recommendedPlan) {
			return { text: "Not enough RAM", color: "red" as const };
		}

		return undefined;
	};

	const currentPlan = {
		...pricingPlans[currentPlanIndex],
		badge: getCurrentPlanBadge(),
	};

	return (
		<form
			onSubmit={handleSubmit(onSubmit, onError)}
			className="flex flex-col items-center justify-center min-h-150 space-y-8 px-4 w-full lg:w-2/3 lg:px-0"
		>
			<div className="flex w-full justify-start">
				<h1 className=" text-2xl font-bold">Create your server</h1>
			</div>

			<div className="flex flex-col sm:flex-row w-full space-y-12 sm:space-x-8">
				<div className="flex flex-col w-full sm:w-4/5 space-y-8">
					<Card className="w-full">
						<CardContent className="space-y-5">
							<h3 className="text-lg">Configure your installation</h3>
							<div className="space-y-3">
								<h4>What should we install on you server?</h4>
								<Select
									onValueChange={(value) => setValue("serverType", value)}
								>
									<SelectTrigger className="w-full my-2">
										<SelectValue placeholder="Select server type" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="vanilla">Vanilla</SelectItem>
										<SelectItem value="paper">Paper</SelectItem>
										<SelectItem value="curseforge">CurseForge</SelectItem>
									</SelectContent>
								</Select>

								{errors.serverType && (
									<p className="text-red-500">
										{errors.serverType.message as string}
									</p>
								)}

								<Select
									onValueChange={(value) => setValue("serverVersion", value)}
								>
									<SelectTrigger className="w-full my-2">
										<SelectValue placeholder="Select server version" />
									</SelectTrigger>
									<SelectContent>
										<SelectItem value="1.21.11">1.21.11</SelectItem>
										<SelectItem value="1.21.10">1.21.10</SelectItem>
										<SelectItem value="1.21.9">1.21.9</SelectItem>
									</SelectContent>
								</Select>
								{errors.serverVersion && (
									<p className="text-red-500">
										{errors.serverVersion.message as string}
									</p>
								)}
							</div>
							<div className="space-y-3">
								<h4>
									How many players do you expect to be on the server at once?
								</h4>
								<Slider
									className="my-2"
									value={watch("playerCount") ? [watch("playerCount")] : [0]}
									onValueChange={(value) => setValue("playerCount", value[0])}
									defaultValue={[0]}
									max={100}
									step={1}
								/>
								<p className="text-sm text-stone-400">
									{watch("playerCount")} players
								</p>
								{errors.playerCount && (
									<p className="text-red-500">
										{errors.playerCount.message as string}
									</p>
								)}
							</div>
						</CardContent>
					</Card>

					<Card className="w-full">
						<CardContent className="space-y-3">
							<h3 className="text-lg">Select your location</h3>
							<Select onValueChange={(value) => setValue("region", value)}>
								<SelectTrigger className="w-full my-2">
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
							{errors.region && (
								<p className="text-red-500">
									{errors.region.message as string}
								</p>
							)}
						</CardContent>
					</Card>

					<Card className="w-full pb-8">
						<CardContent className="space-y-4">
							<h3 className="text-lg">Add your initial configuration</h3>
							<div className="space-y-3">
								<h4>Server name</h4>
								<Input
									onChange={(e) => setValue("serverName", e.target.value)}
								/>
								{errors.serverName && (
									<p className="text-red-500">
										{errors.serverName.message as string}
									</p>
								)}
							</div>

							<div className="flex justify-between space-x-6">
								<div className="w-1/2 space-y-3">
									<h4>Select difficulty</h4>
									<Select
										onValueChange={(value) => setValue("difficulty", value)}
									>
										<SelectTrigger className="w-full">
											<SelectValue placeholder="Select a difficulty" />
										</SelectTrigger>
										<SelectContent>
											<SelectItem value="peaceful">Peaceful</SelectItem>
											<SelectItem value="easy">Easy</SelectItem>
											<SelectItem value="normal">Normal</SelectItem>
											<SelectItem value="hard">Hard</SelectItem>
											<SelectItem value="hardcore">Hardcore</SelectItem>
										</SelectContent>
									</Select>
									{errors.difficulty && (
										<p className="text-red-500">
											{errors.difficulty.message as string}
										</p>
									)}
								</div>
								<div className="w-1/2 flex flex-col justify-center space-y-3 md:space-y-1">
									<h4 className="pb-2">Allow only premium players</h4>
									<Switch
										checked={watch("premiumOnly") || false}
										onCheckedChange={(checked) =>
											setValue("premiumOnly", checked)
										}
									/>
								</div>
							</div>
						</CardContent>
					</Card>

					<button
						type="submit"
						className="w-full bg-yellow-500 hover:bg-yellow-600 text-stone-900 font-semibold py-3 rounded-lg transition-colors"
					>
						Create Server
					</button>
				</div>

				<div className="flex flex-col w-full sm:w-2/5 space-y-8">
					<div className="max-w-2xl w-full overflow-visible">
						<div className="space-y-4 py-8">
							{isLoadingRecommendation && (
								<p className="text-center text-sm text-stone-400">
									Finding best plan...
								</p>
							)}
							{/* Carousel */}
							<div className="relative flex items-center justify-center gap-4">
								{/* Previous Button */}
								<button
									type="button"
									onClick={handlePrevious}
									className="z-10 flex items-center justify-center w-10 h-10 rounded-full bg-stone-800 hover:bg-stone-700 border-2 border-stone-600 transition-colors shrink-0"
									aria-label="Previous plan"
								>
									<HugeiconsIcon icon={ArrowLeft} size={20} />
								</button>

								{/* Pricing Card */}
								<div className="flex-1">
									<PricingCard
										ram={currentPlan.ram}
										price={currentPlan.price}
										badge={currentPlan.badge}
									/>
								</div>

								{/* Next Button */}
								<button
									type="button"
									onClick={handleNext}
									className="z-10 flex items-center justify-center w-10 h-10 rounded-full bg-stone-800 hover:bg-stone-700 border-2 border-stone-600 transition-colors shrink-0"
									aria-label="Next plan"
								>
									<HugeiconsIcon icon={ArrowRight} size={20} />
								</button>
							</div>

							{/* Indicator Dots */}
							<div className="flex justify-center gap-2">
								{pricingPlans.map((_, index) => (
									<button
										key={index}
										type="button"
										onClick={() => {
											setCurrentPlanIndex(index);
											setValue("ramPlan", pricingPlans[index].ram);
										}}
										className={`w-2 h-2 rounded-full transition-all ${
											index === currentPlanIndex
												? "bg-yellow-500 w-8"
												: "bg-stone-600 hover:bg-stone-500"
										}`}
										aria-label={`Go to plan ${index + 1}`}
									/>
								))}
							</div>
						</div>
					</div>
				</div>
			</div>
		</form>
	);
}
