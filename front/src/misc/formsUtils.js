
export function generatePassword() {
	let pwd = "";
	const possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

	for (var i = 0; i < 5; i++) {
		pwd += possible.charAt(Math.floor(Math.random() * possible.length));
	}

	return pwd;
}
