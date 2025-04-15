<script lang="ts">
	import { onMount } from 'svelte';

	let healthStatus = 'Checking...';

	onMount(async () => {
		try {
			const res = await fetch('http://localhost:8081/health');
			if (!res.ok) throw new Error('Not OK');

			const data = await res.json();
			console.log('Health Check:', data);

			healthStatus = data.message || 'OK';
		} catch (err) {
			console.error('Health check failed:', err);
			healthStatus = 'Failed to reach server';
		}
	});
</script>

<p>Server Health: {healthStatus}</p>
