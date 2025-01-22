<script>
	import { preventDefault } from 'svelte/legacy';

	import { onMount } from 'svelte';
	import BarraLateral from '$lib/components/BarraLateral.svelte';

	let users = $state([]);
	let newUser = $state({ username: '', password: '' });
	let changePassword = $state({ username: '', newPassword: '' });
	let error = $state('');
	let success = $state('');

	async function loadUsers() {
		const user = JSON.parse(localStorage.getItem('user'));
		const response = await fetch('http://localhost:5000/admin/usuarios', {
			headers: {
				Authorization: user.token
			}
		});
		if (response.ok) {
			users = await response.json();
		}
	}

	async function createUser() {
		const user = JSON.parse(localStorage.getItem('user'));
		const response = await fetch('http://localhost:5000/admin/criar_usuario', {
			method: 'POST',
			headers: {
				Authorization: user.token,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ ...newUser, role: 'funcionario' })
		});

		if (response.ok) {
			success = 'Usuário criado com sucesso';
			newUser = { username: '', password: '' };
			await loadUsers();
		} else {
			const data = await response.json();
			error = data.error || 'Erro ao criar usuário';
		}
	}

	async function updatePassword() {
		const user = JSON.parse(localStorage.getItem('user'));
		const response = await fetch('http://localhost:5000/admin/alterar_senha', {
			method: 'POST',
			headers: {
				Authorization: user.token,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(changePassword)
		});

		if (response.ok) {
			success = 'Senha alterada com sucesso';
			changePassword = { username: '', newPassword: '' };
		} else {
			const data = await response.json();
			error = data.error || 'Erro ao alterar senha';
		}
	}

	onMount(loadUsers);
</script>

<div class="page">
	<BarraLateral />
	<div class="content">
		<h1>Gerenciamento de Usuários</h1>

		{#if error}
			<div class="error">{error}</div>
		{/if}
		{#if success}
			<div class="success">{success}</div>
		{/if}

		<div class="section">
			<h2>Criar Novo Usuário</h2>
			<form onsubmit={preventDefault(createUser)}>
				<div class="input-group">
					<label for="username">Usuário:</label>
					<input type="text" id="username" bind:value={newUser.username} required />
				</div>
				<div class="input-group">
					<label for="password">Senha:</label>
					<input type="password" id="password" bind:value={newUser.password} required />
				</div>
				<button type="submit">Criar Usuário</button>
			</form>
		</div>

		<div class="section">
			<h2>Alterar Senha</h2>
			<form onsubmit={preventDefault(updatePassword)}>
				<div class="input-group">
					<label for="change-username">Usuário:</label>
					<input type="text" id="change-username" bind:value={changePassword.username} required />
				</div>
				<div class="input-group">
					<label for="new-password">Nova Senha:</label>
					<input
						type="password"
						id="new-password"
						bind:value={changePassword.newPassword}
						required
					/>
				</div>
				<button type="submit">Alterar Senha</button>
			</form>
		</div>

		<div class="section">
			<h2>Usuários Cadastrados</h2>
			<table>
				<thead>
					<tr>
						<th>Usuário</th>
						<th>Função</th>
					</tr>
				</thead>
				<tbody>
					{#each users as user}
						<tr>
							<td>{user.username}</td>
							<td>{user.role}</td>
						</tr>
					{/each}
				</tbody>
			</table>
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

	.section {
		margin: 20px 0;
		padding: 20px;
		background: white;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	}

	.input-group {
		margin-bottom: 15px;
	}

	label {
		display: block;
		margin-bottom: 5px;
	}

	input {
		width: 100%;
		padding: 8px;
		border: 1px solid #ddd;
		border-radius: 4px;
	}

	button {
		background: #4caf50;
		color: white;
		padding: 10px 20px;
		border: none;
		border-radius: 4px;
		cursor: pointer;
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

	.error {
		color: red;
		margin: 10px 0;
	}

	.success {
		color: green;
		margin: 10px 0;
	}
</style>
