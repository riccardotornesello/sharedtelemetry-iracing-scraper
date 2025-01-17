<script lang="ts">
	import type { PageData } from './$types';
	import dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';
	import TimeCell from '../components/time-cell.svelte';

	dayjs.extend(duration);

	function formatSeconds(seconds: number): string {
		const totalMilliseconds = Math.floor(seconds * 1000); // Convert seconds to milliseconds
		const durationObj = dayjs.duration(totalMilliseconds);
		return durationObj.format('mm:ss.SSS');
	}

	function formatDate(date: string): string {
		return dayjs(date).format('DD/MM');
	}

	let { data }: { data: PageData } = $props();
</script>

<table class="w-full table-auto text-left text-sm text-gray-500 rtl:text-right dark:text-gray-400">
	<thead
		class="bg-gray-50 text-center text-xs uppercase text-gray-700 dark:bg-gray-700 dark:text-gray-400"
	>
		<tr>
			<th scope="col" class="px-6 py-3" rowspan="2" colspan="2">Pilota</th>
			<th scope="col" class="px-6 py-3" rowspan="2">Somma</th>
			{#each Object.entries(data.dates) as [track, dates]}
				<th scope="col" class="px-6 py-3" colspan={dates.length}>{track}</th>
			{/each}
		</tr>
		<tr>
			{#each Object.entries(data.dates) as [track, dates]}
				{#each dates as date}
					<th scope="col" class="px-6 py-3">{formatDate(date)}</th>
				{/each}
			{/each}
		</tr>
	</thead>
	<tbody>
		{#each data.results as result, index}
			<tr class="border-b bg-white dark:border-gray-700 dark:bg-gray-800">
				<td class="px-6 py-4 text-center">P{index + 1}</td>
				<td class="px-6 py-4">{result.name || result.custId}</td>
				<td class="px-6 py-4 font-mono text-center">{formatSeconds(result.sum)}</td>
				{#each Object.entries(data.dates) as [track, dates]}
					{#each dates as date}
						<td class="px-2 py-2">
							<TimeCell {result} {date} {track} bestPerTrack={data.bestPerTrack} />
						</td>
					{/each}
				{/each}
			</tr>
		{/each}
	</tbody>
</table>
