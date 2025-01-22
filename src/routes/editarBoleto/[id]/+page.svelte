<script>
	import { preventDefault } from 'svelte/legacy';

	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import BarraLateral from '$lib/components/BarraLateral.svelte';

	let boleto = $state({
		nome: '',
		dataVencimento: '',
		valor: '',
		pdfPath: ''
	});

	let id = $page.params.id;

	onMount(async () => {
		try {
			const response = await fetch(`http://localhost:5000/buscar_boleto/${id}`);
			if (response.ok) {
				const data = await response.json();
				boleto = data.boleto;
			}
		} catch (error) {
			console.error('Erro ao carregar boleto:', error);
		}
	});

	async function handleSubmit(event) {
		event.preventDefault();
		try {
			const response = await fetch(`http://localhost:5000/atualizar_boleto/${id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					...boleto,
					valor: parseFloat(boleto.valor)
				})
			});

			if (response.ok) {
				alert('Boleto atualizado com sucesso!');
				goto('/boletos');
			} else {
				throw new Error('Erro ao atualizar boleto');
			}
		} catch (error) {
			alert('Erro ao atualizar boleto: ' + error.message);
		}
	}

	async function handleFileUpload(event) {
		const file = event.target.files[0];
		if (file && file.type === 'application/pdf') {
			const formData = new FormData();
			formData.append('pdf', file);

			try {
				const response = await fetch('http://localhost:5000/upload_boleto', {
					method: 'POST',
					body: formData
				});

				if (response.ok) {
					const data = await response.json();
					boleto.pdfPath = data.filepath;
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
</script>

<div class="container">
	<BarraLateral />
	<div class="content">
		<div class="form-container">
			<h1>Editar Boleto</h1>
			<form onsubmit={preventDefault(handleSubmit)}>
				<div class="form-group">
					<label for="nome">Nome:</label>
					<input type="text" id="nome" bind:value={boleto.nome} required />
				</div>

				<div class="form-group">
					<label for="dataVencimento">Data de Vencimento:</label>
					<input type="date" id="dataVencimento" bind:value={boleto.dataVencimento} required />
				</div>

				<div class="form-group">
					<label for="valor">Valor:</label>
					<input type="number" id="valor" bind:value={boleto.valor} step="0.01" required />
				</div>

				<div class="form-group">
					<label for="pdfFile">PDF do Boleto:</label>
					<input type="file" id="pdfFile" accept=".pdf" onchange={handleFileUpload} />
					{#if boleto.pdfPath}
						<span class="file-name">Arquivo atual: {boleto.pdfPath}</span>
					{/if}
				</div>

				<div class="button-group">
					<button type="submit" class="save-btn">Salvar Alterações</button>
					<button type="button" class="cancel-btn" onclick={() => goto('/boletos')}>
						Cancelar
					</button>
				</div>
			</form>
		</div>
	</div>
</div>

<style>
	/* Mesmo estilo do registrarMaterial */
	/* ...existing styles... */
</style>
