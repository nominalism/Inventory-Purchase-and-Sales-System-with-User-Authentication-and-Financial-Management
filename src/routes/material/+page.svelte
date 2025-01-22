<script>
	import BarraLateral from '$lib/components/BarraLateral.svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation'; // Para redirecionar ao editar

	let materiais = $state([]);
	let busca = $state('');
	let error = '';
	let user = null;

	onMount(() => {
		const storedUser = JSON.parse(localStorage.getItem('user'));
		if (storedUser) {
			user = storedUser;
		}
		buscarMateriais();
	});

	async function buscarMateriais() {
		try {
			const response = await fetch(`http://localhost:5000/busca_material?busca=${busca}`, {
				headers: {
					Authorization: JSON.parse(localStorage.getItem('user')).token
				}
			});

			if (response.ok) {
				const data = await response.json();
				materiais = data.materiais;
			}
		} catch (err) {
			error = 'Erro ao conectar com o servidor';
			console.error(err);
		}
	}

	// Função para excluir material
	async function excluirMaterial(material) {
		if (confirm('Tem certeza que deseja excluir este material?')) {
			try {
				const user = JSON.parse(localStorage.getItem('user'));
				console.log('Excluindo material ID:', material.id); // Changed from ID to id
				const response = await fetch(
					`http://localhost:5000/admin/excluir_material/${material.id}`,
					{
						method: 'DELETE',
						headers: {
							Authorization: user.token
						}
					}
				);

				const data = await response.json();

				if (response.ok) {
					alert('Material excluído com sucesso!');
					await buscarMateriais();
				} else {
					throw new Error(data.error || 'Erro ao excluir material');
				}
			} catch (error) {
				console.error('Erro ao excluir material:', error);
				alert(error.message || 'Erro ao excluir material');
			}
		}
	}

	// Função para editar material
	function editarMaterial(codigo) {
		goto(`/editar/${codigo}`); // Redireciona para a página de edição
	}

	async function comprarMaterial(material) {
		const quantidade = parseInt(prompt('Quantidade a comprar:'));
		if (!quantidade || quantidade <= 0) return;

		try {
			const user = JSON.parse(localStorage.getItem('user'));
			console.log('Material ID:', material.id); // Debug log
			const response = await fetch(`http://localhost:5000/admin/comprar_material/${material.id}`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: user.token
				},
				body: JSON.stringify({ quantidade })
			});

			const data = await response.json();
			if (!response.ok) {
				throw new Error(data.error || 'Erro ao comprar material');
			}

			await buscarMateriais();
			alert('Compra realizada com sucesso!');
		} catch (error) {
			console.error('Erro:', error);
			alert(error.message || 'Erro ao conectar com o servidor');
		}
	}

	async function venderMaterial(material) {
		const quantidade = parseInt(prompt('Quantidade a vender:'));
		if (!quantidade || quantidade <= 0) return;

		try {
			const user = JSON.parse(localStorage.getItem('user'));
			console.log('Material ID:', material.id); // Debug log
			const response = await fetch(`http://localhost:5000/vender_material/${material.id}`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: user.token
				},
				body: JSON.stringify({ quantidade })
			});

			const data = await response.json();
			if (!response.ok) {
				throw new Error(data.error || 'Erro ao vender material');
			}

			await buscarMateriais();
			alert('Venda realizada com sucesso!');
		} catch (error) {
			console.error('Erro:', error);
			alert(error.message || 'Erro ao conectar com o servidor');
		}
	}

	// Função para formatar data
	function formatarData(data) {
		if (!data) return '';
		const [ano, mes, dia] = data.split('-');
		return `${dia.slice(0, 2)}/${mes}/${ano}`;
	}
</script>

<div class="container">
	<BarraLateral></BarraLateral>

	<div class="conteudo">
		<div class="pesquisa">
			<input
				type="text"
				placeholder="Busque pelo código ou nome do material"
				class="pesquisaMaterial"
				bind:value={busca}
				oninput={buscarMateriais}
			/>
		</div>

		<table class="tabelaMateriais">
			<thead>
				<tr>
					<th>Código</th>
					<th>Nome</th>
					<th>Valor de Venda</th>
					<th>Estoque</th>
					<th>Local</th>
					<th>Data de Compra</th>
					<th>Ações</th>
				</tr>
			</thead>
			<tbody>
				{#if materiais.length > 0}
					{#each materiais as material}
						<tr>
							<td>{material.codigo}</td>
							<td>{material.nome}</td>
							<td>R$ {material.valorVenda}</td>
							<td>{material.estoque}</td>
							<td>{material.local}</td>
							<td>{formatarData(material.dataCompra)}</td>
							<td>
								<button onclick={() => editarMaterial(material.id)}>Editar</button>
								<button onclick={() => excluirMaterial(material)}>Excluir</button>
								<button onclick={() => comprarMaterial(material)}>Comprar</button>
								<button onclick={() => venderMaterial(material)}>Vender</button>
							</td>
						</tr>
					{/each}
				{:else}
					<tr>
						<td colspan="7" style="text-align: center;">Nenhum material encontrado.</td>
					</tr>
				{/if}
			</tbody>
		</table>
	</div>
</div>

<style>
	/* Layout principal */
	.container {
		display: flex;
		height: 100vh;
		margin: 0;
	}

	/* Conteúdo à direita */
	.conteudo {
		flex: 1;
		padding: 20px;
		overflow-y: auto;
	}

	/* Estilo da pesquisa */
	.pesquisa {
		margin-bottom: 20px;
	}
	.pesquisaMaterial {
		width: 100%;
		padding: 10px;
		font-size: 16px;
	}

	/* Tabela */
	.tabelaMateriais {
		width: 100%;
		border-collapse: collapse;
	}
	.tabelaMateriais th,
	.tabelaMateriais td {
		border: 1px solid #ddd;
		padding: 8px;
		text-align: left;
	}
	.tabelaMateriais th {
		background-color: #f4f4f4;
	}
	.tabelaMateriais button {
		margin: 0 5px;
		padding: 5px 10px;
		font-size: 14px;
		cursor: pointer;
	}
</style>
