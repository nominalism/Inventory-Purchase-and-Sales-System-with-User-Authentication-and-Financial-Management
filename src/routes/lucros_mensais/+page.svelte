<script>
	import { onMount } from 'svelte';
	import BarraLateral from '$lib/components/BarraLateral.svelte';

	let lucros = $state({
		meses: [],
		valores: []
	});

	onMount(async () => {
		try {
			const user = JSON.parse(localStorage.getItem('user'));
			const response = await fetch('http://localhost:5000/lucros_mensais', {
				headers: {
					Authorization: user.token
				}
			});
			if (response.ok) {
				lucros = await response.json();
			}
		} catch (error) {
			console.error('Erro ao carregar lucros:', error);
		}
	});
</script>

<div class="page">
	<BarraLateral />
	<div class="content">
		<div class="container">
			<h1>Lucros Mensais</h1>

			<div class="lucros-table">
				<table>
					<thead>
						<tr>
							<th>Mês</th>
							<th>Valor</th>
						</tr>
					</thead>
					<tbody>
						{#each lucros.meses || [] as mes, i}
							<tr>
								<td>{mes}</td>
								<td>R$ {lucros.valores[i].toFixed(2)}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</div>

<style>
	.page {
		display: flex;
		min-height: 100vh;
	}

	.content {
		flex: 1;
		padding: 20px;
	}

	.container {
		padding: 20px;
		max-width: 800px;
		margin: 0 auto;
	}

	h1 {
		color: #333;
		margin-bottom: 20px;
	}

	.lucros-table {
		width: 100%;
		overflow-x: auto;
	}

	table {
		width: 100%;
		border-collapse: collapse;
		margin-top: 20px;
	}

	th,
	td {
		padding: 12px;
		text-align: left;
		border-bottom: 1px solid #ddd;
	}

	th {
		background-color: #f5f5f5;
		font-weight: bold;
	}

	tr:hover {
		background-color: #f9f9f9;
	}
</style>
