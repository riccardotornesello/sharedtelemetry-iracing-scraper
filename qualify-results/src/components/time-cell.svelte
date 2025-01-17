<script lang="ts">
	import dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';

	interface Props {
		result: any;
		date: string;
		track: string;
		bestPerTrack: any;
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
	<div
		class="h-full w-full rounded p-4 font-mono text-gray-300"
		class:bg-purple-500={isOverallBest}
		class:text-purple-300={isOverallBest}
		class:bg-green-500={isPersonalBest && !isOverallBest}
		class:text-green-300={isPersonalBest && !isOverallBest}
	>
		{formatSeconds(result.laps[date].time)}
	</div>
{:else}
	<div class="h-full w-full rounded p-4 font-mono text-gray-300">-</div>
{/if}
