<script lang="ts">
	import dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';
	import TimeCard from './time-card.svelte';
	import type { DriverResult } from '../routes/types';

	interface Props {
		result: DriverResult;
		date: string;
		track: string;
		bestPerTrack: Record<string, number>;
	}

	let { result, date, track, bestPerTrack }: Props = $props();

	dayjs.extend(duration);

	function formatSeconds(seconds: number): string {
		const totalMilliseconds = Math.floor(seconds * 1000); // Convert seconds to milliseconds
		const durationObj = dayjs.duration(totalMilliseconds);
		return durationObj.format('mm:ss.SSS');
	}

	let isOverallBest = result.laps[date] && result.laps[date].time === bestPerTrack[track];
	let isPersonalBest = result.laps[date] && result.laps[date].isBest === true;
</script>

{#if result.laps[date]}
	<TimeCard {isPersonalBest} {isOverallBest} time={result.laps[date].time} />
{:else}
	<TimeCard />
{/if}
