<script>
	import { preventDefault } from 'svelte/legacy';

	import { onMount } from 'svelte';
	import BarraLateral from '$lib/components/BarraLateral.svelte';

	let boletos = $state([]);
	let novoBoletoDados = $state({
		nome: '',
		dataVencimento: '',
		valor: '',
		pdfPath: ''
	});

	let totalBoletos = $state(0);

	function calcularTotal(boletos) {
		return boletos.reduce((sum, boleto) => sum + parseFloat(boleto.valor), 0);
	}

	async function carregarBoletos() {
		try {
			const user = JSON.parse(localStorage.getItem('user'));
			const response = await fetch('http://localhost:5000/buscar_boletos', {
				headers: {
					Authorization: user.token
				}
			});
			if (!response.ok) {
				throw new Error('Erro ao buscar boletos');
			}
			const data = await response.json();
			if (data && data.boletos) {
				boletos = data.boletos;
				totalBoletos = calcularTotal(boletos);
			}
		} catch (error) {
			console.error('Erro ao carregar boletos:', error);
			boletos = [];
		}
	}

	onMount(() => {
		carregarBoletos();
	});

	async function registrarBoleto(event) {
		event.preventDefault();
		try {
			const user = JSON.parse(localStorage.getItem('user'));
			const boletoParaEnviar = {
				...novoBoletoDados,
				valor: parseFloat(novoBoletoDados.valor) // Converte para número
			};

			const response = await fetch('http://localhost:5000/admin/adicionar_boleto', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: user.token
				},
				body: JSON.stringify(boletoParaEnviar)
			});

			if (!response.ok) {
				throw new Error('Erro ao adicionar boleto');
			}

			alert('Boleto registrado com sucesso!');
			await carregarBoletos();

			// Limpar formulário
			novoBoletoDados = {
				nome: '',
				dataVencimento: '',
				valor: '',
				pdfPath: ''
			};
		} catch (error) {
			console.error('Erro ao registrar boleto:', error);
			alert('Erro ao registrar boleto');
		}
	}

	async function excluirBoleto(id) {
		if (!confirm('Tem certeza que deseja excluir este boleto?')) return;

		try {
			const user = JSON.parse(localStorage.getItem('user'));
			const response = await fetch(`http://localhost:5000/admin/excluir_boleto/${id}`, {
				method: 'DELETE',
				headers: {
					Authorization: user.token
				}
			});

			if (!response.ok) {
				throw new Error('Erro ao excluir boleto');
			}

			alert('Boleto excluído com sucesso!');
			await carregarBoletos();
		} catch (error) {
			console.error('Erro ao excluir boleto:', error);
			alert('Erro ao excluir boleto');
		}
	}

	function formatarData(data) {
		if (!data) return '';
		const date = new Date(data);
		return date.toLocaleDateString('pt-BR');
	}

	async function handleFileUpload(event) {
		const file = event.target.files[0];
		if (file && file.type === 'application/pdf') {
			const formData = new FormData();
			formData.append('pdf', file);

			try {
				const user = JSON.parse(localStorage.getItem('user'));
				const response = await fetch('http://localhost:5000/admin/upload_boleto', {
					method: 'POST',
					headers: {
						Authorization: user.token
					},
					body: formData
				});

				if (response.ok) {
					const data = await response.json();
					novoBoletoDados.pdfPath = data.filepath;
				} else {
					alert('Erro ao fazer upload do arquivo');
				}
			} catch (error) {
				console.error('Erro no upload:', error);
				alert('Erro ao fazer upload do arquivo');
			}
		} else {
			alert('Por favor, selecione um arquivo PDF');
		}
	}

	async function pagarBoleto(boleto) {
		if (
			!confirm(
				`Tem certeza que deseja pagar o boleto ${boleto.nome} no valor de R$ ${parseFloat(boleto.valor).toFixed(2)}?`
			)
		) {
			return;
		}

		try {
			const user = JSON.parse(localStorage.getItem('user'));
			const response = await fetch(`http://localhost:5000/admin/pagar_boleto/${boleto.id}`, {
				method: 'POST',
				headers: {
					Authorization: user.token
				}
			});

			if (!response.ok) {
				const error = await response.json();
				throw new Error(error.error || 'Erro ao pagar boleto');
			}

			alert('Boleto pago com sucesso!');
			await carregarBoletos();
		} catch (error) {
			alert(error.message);
		}
	}

	function editarBoleto(boleto) {
		goto(`/editarBoleto/${boleto.id}`);
	}
</script>

<div class="container">
	<BarraLateral />
	<div class="content">
		<div class="header">
			<h1>Gerenciamento de Boletos - Valor Total: R$ {totalBoletos.toFixed(2)}</h1>
		</div>

		<div class="form-section">
			<h2>Registrar Novo Boleto</h2>
			<form onsubmit={preventDefault(registrarBoleto)}>
				<div class="input-group">
					<label for="nome">Nome:</label>
					<input type="text" id="nome" bind:value={novoBoletoDados.nome} required />
				</div>

				<div class="input-group">
					<label for="dataVencimento">Data de Vencimento:</label>
					<input
						type="date"
						id="dataVencimento"
						bind:value={novoBoletoDados.dataVencimento}
						required
					/>
				</div>

				<div class="input-group">
					<label for="valor">Valor:</label>
					<input type="number" id="valor" bind:value={novoBoletoDados.valor} step="0.01" required />
				</div>

				<div class="input-group">
					<label for="pdfFile">PDF do Boleto:</label>
					<input
						type="file"
						id="pdfFile"
						accept=".pdf"
						onchange={handleFileUpload}
						class="file-input"
					/>
					{#if novoBoletoDados.pdfPath}
						<span class="file-name">Arquivo selecionado</span>
					{/if}
				</div>

				<button type="submit">Registrar Boleto</button>
			</form>
		</div>

		<div class="table-section">
			<h2>Boletos Cadastrados</h2>
			<table>
				<thead>
					<tr>
						<th>Nome</th>
						<th>Vencimento</th>
						<th>Valor</th>
						<th>PDF</th>
						<th>Ações</th>
					</tr>
				</thead>
				<tbody>
					{#if boletos.length === 0}
						<tr>
							<td colspan="5" class="empty-message">Nenhum boleto cadastrado</td>
						</tr>
					{:else}
						{#each boletos as boleto}
							<tr>
								<td>{boleto.nome}</td>
								<td>{formatarData(boleto.dataVencimento)}</td>
								<td>R$ {parseFloat(boleto.valor).toFixed(2)}</td>
								<td>
									{#if boleto.pdfPath}
										<a href={boleto.pdfPath} target="_blank">Ver PDF</a>
									{/if}
								</td>
								<td>
									<div class="btn-group">
										<button class="edit-btn" onclick={() => editarBoleto(boleto)}>Editar</button>
										<button class="pay-btn" onclick={() => pagarBoleto(boleto)}>Pagar</button>
										<button class="delete-btn" onclick={() => excluirBoleto(boleto.id)}
											>Excluir</button
										>
									</div>
								</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
	</div>
</div>

<style>
	.container {
		display: flex;
		height: 100vh;
	}

	.content {
		flex: 1;
		padding: 20px;
		background-color: #f5f5f5;
	}

	.header {
		margin-bottom: 20px;
	}

	.form-section {
		background: white;
		padding: 20px;
		border-radius: 8px;
		margin-bottom: 20px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.input-group {
		margin-bottom: 15px;
	}

	.input-group label {
		display: block;
		margin-bottom: 5px;
	}

	.input-group input {
		width: 100%;
		padding: 8px;
		border: 1px solid #ddd;
		border-radius: 4px;
	}

	table {
		width: 100%;
		background: white;
		border-radius: 8px;
		border-collapse: collapse;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	th,
	td {
		padding: 12px;
		text-align: left;
		border-bottom: 1px solid #ddd;
	}

	th {
		background-color: #f8f9fa;
	}

	.empty-message {
		text-align: center;
		padding: 20px;
		color: #666;
	}

	button {
		padding: 8px 16px;
		background-color: #4caf50;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}

	button:hover {
		background-color: #45a049;
	}

	.delete-btn {
		background-color: #dc3545;
	}

	.delete-btn:hover {
		background-color: #c82333;
	}

	a {
		color: #007bff;
		text-decoration: none;
	}

	a:hover {
		text-decoration: underline;
	}

	.file-input {
		padding: 8px;
		border: 1px solid #ddd;
		border-radius: 4px;
		width: 100%;
	}

	.file-name {
		color: #4caf50;
		font-size: 0.9em;
		margin-top: 4px;
		display: block;
	}

	.btn-group {
		display: flex;
		gap: 5px;
	}

	.edit-btn {
		background-color: #4a90e2;
	}

	.edit-btn:hover {
		background-color: #357abd;
	}

	.pay-btn {
		background-color: #28a745;
	}

	.pay-btn:hover {
		background-color: #218838;
	}
</style>
