<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(() => {
		// Check if user is logged in
		const user = localStorage.getItem('user');
		if (!user) {
			goto('/login');
		} else {
			goto('/material');
		}
	});
</script>

<div>Redirecionando...</div>
