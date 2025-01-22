<script>
	import { onMount } from 'svelte';
	import { writable } from 'svelte/store';

	let user = writable(null);

	onMount(() => {
		const storedUser = localStorage.getItem('user');
		if (storedUser) {
			$user = JSON.parse(storedUser);
		}
	});

	function handleLogout() {
		localStorage.removeItem('user');
		window.location.href = '/login';
	}
</script>

<nav class="sidebar">
	<ul>
		<li><a href="/material">Materiais</a></li>
		{#if $user?.user?.role === 'admin'}
			<li><a href="/material/registrar">Registrar Material</a></li>

			<li><a href="/boletos">Boletos</a></li>
			<li><a href="/financeiro">Caixa</a></li>
			<li><a href="/lucros_mensais">Lucros Mensais</a></li>
			<li><a href="/usuarios">Usuários</a></li>
		{/if}
		<li><button onclick={handleLogout}>Sair</button></li>
	</ul>
</nav>

<style>
	.sidebar {
		width: 200px;
		background-color: #333;
		height: 100vh;
		padding: 20px;
	}

	ul {
		list-style: none;
		padding: 0;
	}

	li {
		margin-bottom: 10px;
	}

	a,
	button {
		color: white;
		text-decoration: none;
		display: block;
		padding: 10px;
		border-radius: 4px;
		width: 100%;
		text-align: left;
		background: none;
		border: none;
		font-size: 1em;
		cursor: pointer;
	}

	a:hover,
	button:hover {
		background-color: #444;
	}
</style>
