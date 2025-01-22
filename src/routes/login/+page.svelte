<script>
	import { goto } from '$app/navigation';
	let username = $state('');
	let password = $state('');
	let error = $state('');

	async function handleLogin(event) {
		event.preventDefault();
		try {
			const response = await fetch('http://localhost:5000/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ username, password })
			});

			if (response.ok) {
				const data = await response.json();
				localStorage.setItem('user', JSON.stringify(data));
				goto('/material');
			} else {
				error = 'Usuário ou senha inválidos';
			}
		} catch (err) {
			error = 'Erro ao conectar com o servidor';
		}
	}
</script>

<div class="login-container">
	<div class="login-box">
		<h1>Login</h1>
		{#if error}
			<div class="error">{error}</div>
		{/if}
		<form onsubmit={handleLogin}>
			<div class="input-group">
				<label for="username">Usuário:</label>
				<input type="text" id="username" bind:value={username} required />
			</div>
			<div class="input-group">
				<label for="password">Senha:</label>
				<input type="password" id="password" bind:value={password} required />
			</div>
			<button type="submit">Entrar</button>
		</form>
	</div>
</div>

<style>
	.login-container {
		display: flex;
		justify-content: center;
		align-items: center;
		height: 100vh;
		background-color: #f5f5f5;
	}

	.login-box {
		background: white;
		padding: 30px;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
		width: 100%;
		max-width: 400px;
	}

	.input-group {
		margin-bottom: 20px;
	}

	label {
		display: block;
		margin-bottom: 8px;
	}

	input {
		width: 100%;
		padding: 8px;
		border: 1px solid #ddd;
		border-radius: 4px;
	}

	button {
		width: 100%;
		padding: 10px;
		background-color: #4caf50;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}

	.error {
		color: red;
		margin-bottom: 15px;
	}
</style>
