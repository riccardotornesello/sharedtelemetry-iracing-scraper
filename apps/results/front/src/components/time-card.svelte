<script lang="ts">
	import dayjs from 'dayjs';
	import duration from 'dayjs/plugin/duration';

	interface Props {
		time?: number;
		isOverallBest?: boolean;
		isPersonalBest?: boolean;
	}

	let { time, isOverallBest, isPersonalBest }: Props = $props();

	dayjs.extend(duration);

	function formatMilliseconds(totalMilliseconds: number): string {
		const durationObj = dayjs.duration(totalMilliseconds);
		return durationObj.format('mm:ss.SSS');
	}
</script>

<div
	class="h-full w-full rounded p-3 text-center font-mono text-gray-300"
	class:bg-purple-600={isOverallBest}
	class:text-purple-200={isOverallBest}
	class:bg-green-600={isPersonalBest && !isOverallBest}
	class:text-green-200={isPersonalBest && !isOverallBest}
>
	{time ? formatMilliseconds(time) : '-'}
</div>
