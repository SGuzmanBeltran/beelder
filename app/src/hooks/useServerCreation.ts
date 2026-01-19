import { useCallback, useEffect, useState } from "react";

import axios from "axios";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

const API_URL = import.meta.env.VITE_API_URL;

export const pricingPlans = [
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
	onlineMode: z.boolean(),
});

interface CreationResponse {
	message: string;
	name: string;
	id: string;
}

type ServerConfig = z.infer<typeof serverConfigSchema>;

export function useServerCreation() {
	const [serverVersions, setServerVersions] = useState<string[]>([]);
	const [loadingServerVersions, setLoadingServerVersions] = useState(false);
	const [currentPlanIndex, setCurrentPlanIndex] = useState(0);
	const [isLoadingRecommendation, setIsLoadingRecommendation] = useState(false);
	const [recommendedPlan, setRecommendedPlan] = useState<number | null>(null);
	const [isSubmitting, setIsSubmitting] = useState(false);

	const form = useForm<ServerConfig>({
		resolver: zodResolver(serverConfigSchema),
		defaultValues: {
			playerCount: 0,
			onlineMode: true,
			ramPlan: pricingPlans[0].ram,
		},
	});

	const handlePrevious = () => {
		setCurrentPlanIndex((prev) => {
			const newIndex = prev > 0 ? prev - 1 : pricingPlans.length - 1;
			form.setValue("ramPlan", pricingPlans[newIndex].ram);
			return newIndex;
		});
	};

	const handleNext = () => {
		setCurrentPlanIndex((prev) => {
			const newIndex = prev < pricingPlans.length - 1 ? prev + 1 : 0;
			form.setValue("ramPlan", pricingPlans[newIndex].ram);
			return newIndex;
		});
	};

	// Get recommended plan from backend
	const fetchRecommendedPlan = useCallback(
		async (serverType: string, playerCount: number, region: string) => {
			if (!serverType || !playerCount || !region) return;

			setIsLoadingRecommendation(true);

			let retries = 0;
			const maxRetries = 3;

			while (retries < maxRetries) {
				try {
					const { data } = await axios.get(
						`${API_URL}/api/v1/server/recommended-plans?server_type=${serverType}&player_count=${playerCount}&region=${region}`,
					);

					// Find the index of the recommended plan
					const recommendedRam = data.data.recommendation;
					const planIndex = pricingPlans.findIndex(
						(plan) => plan.ram === recommendedRam,
					);

					if (planIndex !== -1 && planIndex !== 0) {
						// Don't auto-select free plan
						setCurrentPlanIndex(planIndex);
						setRecommendedPlan(planIndex);
						form.setValue("ramPlan", pricingPlans[planIndex].ram);
					}

					// Success - exit the loop
					setIsLoadingRecommendation(false);
					return;
				} catch (error) {
					retries++;

					if (retries === maxRetries) {
						// Failed after all retries
						const message = axios.isAxiosError(error)
							? error.response?.data?.error || error.message
							: "Network error. Please check your connection.";

						toast.error("Failed to fetch recommendation", {
							description: message,
						});
					}

					// Wait before retrying (exponential backoff: 500ms, 1s, 2s)
					if (retries < maxRetries) {
						await new Promise((resolve) =>
							setTimeout(resolve, 500 * Math.pow(2, retries - 1)),
						);
					}
				}
			}

			setIsLoadingRecommendation(false);
		},
		[setIsLoadingRecommendation, setCurrentPlanIndex, setRecommendedPlan, form],
	);

	const fetchServerVersions = useCallback(
		async (serverType: string) => {
			if (!serverType) return;
			setLoadingServerVersions(true);
			form.setValue("serverVersion", "" as string);
			setServerVersions([]);
			let retries = 0;
			const maxRetries = 3;
			while (retries < maxRetries) {
				try {
					const { data } = await axios.get(
						`${API_URL}/api/v1/server/${serverType}/versions`,
					);
					setServerVersions(data.data.versions);
					return;
				} catch (error) {
					retries++;

					if (retries === maxRetries) {
						// Failed after all retries
						const message = axios.isAxiosError(error)
							? error.response?.data?.error || error.message
							: "Network error. Please check your connection.";

						toast.error("Failed to fetch server versions", {
							description: message,
						});
						setServerVersions([]);
					}

					// Wait before retrying (exponential backoff: 500ms, 1s, 2s)
					if (retries < maxRetries) {
						await new Promise((resolve) =>
							setTimeout(resolve, 500 * Math.pow(2, retries - 1)),
						);
					}
				} finally {
					setLoadingServerVersions(false);
				}
			}
		},
		[setServerVersions, form],
	);

	const sendFormData = async (data: ServerConfig) => {
		setIsSubmitting(true);
		try {
			const { data: responseData } = await axios.post<CreationResponse>(
				`${API_URL}/api/v1/server`,
				{
					name: data.serverName,
					server_version: data.serverVersion,
					server_type: data.serverType,
					player_count: data.playerCount,
					region: data.region,
					ram_plan: data.ramPlan,
					difficulty: data.difficulty,
					online_mode: data.onlineMode,
				},
			);
			toast.success("Server created successfully!", {
				description: `Your server "${responseData.name}" has been created.`,
			});
		} catch (error) {
			await new Promise((resolve) => setTimeout(resolve, 1000));
			if (axios.isAxiosError(error)) {
				const message = error.response?.data?.error || error.message;
				toast.error("Failed to create server", {
					description: message,
				});
			} else {
				toast.error("Failed to create server", {
					description: "Network error. Please check your connection.",
				});
			}
		} finally {
			setIsSubmitting(false);
		}
	};

	// Watch for changes in serverType, playerCount, and region
	useEffect(() => {
		const subscription = form.watch((value, { name }) => {
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
						value.region!,
					);
				}, 1000);
				return () => clearTimeout(timeoutId);
			}
		});
		return () => subscription.unsubscribe();
	}, [fetchRecommendedPlan, form]);

	useEffect(() => {
		const subscription = form.watch((value, { name }) => {
			if (name === "serverType" && value.serverType) {
				// Fetch server versions based on server type
				fetchServerVersions(value.serverType!);
				return () => {};
			}
		});
		return () => subscription.unsubscribe();
	}, [fetchServerVersions, form]);

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

	return {
		// Estados
		currentPlanIndex,
		isLoadingRecommendation,
		recommendedPlan,
		isSubmitting,
		serverVersions,
		loadingServerVersions,

		// Form methods
		form,

		// Handlers
		handlePrevious,
		handleNext,
		sendFormData,
		setCurrentPlanIndex,

		// Computed values
		currentPlan: {
			...pricingPlans[currentPlanIndex],
			badge: getCurrentPlanBadge(),
		},
	};
}
