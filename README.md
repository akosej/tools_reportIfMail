# ReportIfMail

ReportIfMail es una aplicación de línea de comandos que monitorea el archivo de registro de Postfix y actualiza una base de datos SQLite3 con información sobre los correos electrónicos enviados a una dirección específica. También proporciona una interfaz web para ver los correos electrónicos recibidos.

## Características

- Monitorea el archivo de registro de Postfix en segundo plano.
- Actualiza una base de datos SQLite3 con información sobre los correos electrónicos enviados a una dirección específica.
- Proporciona una interfaz web para ver los correos electrónicos recibidos.
- Agrupa las líneas de registro por dirección de correo electrónico.

## Instalación

```sh
git clone https://github.com/akosej/ReportIfMail.git
# Descargar dependencias
go mod tidy
go mod vendor
# Compilar
go build -o report main.go
# Copiar a la carpeta que desee alojar el binario
cp report <path>/
# Permiso de ejecución
chmod +x <path>/report
# Crear carpeta necesaria
mkdir <path>/locale

```
### Crear ficheros necesarios
```sh
# Variables necesarias para el sistema
touch 'pathLog = "/var/log/postfix/mail.log"
port="8000"
toEmail="transitoserpen@mail.ho.bpa.cu"' > <path>/locale/.env
# Fichero para crear la lista de los remitentes a filtrar
# El valor de key es el # que va despues de @ y antes del primer (.) para cada agencia 
touch '{
    "6632": "Velasco",
    "key": "...",
    "key": "..."
}' > <path>/locale/agency.json

# Fichero html para renderizar los resultados
touch '<!DOCTYPE html>
<html>
<head>
	<title>Reporte de balance</title>
</head>
<body>
	<h1>Reporte de balance</h1>
	<table>
		<tr>
			<th>ID</th>
			<th>Fecha</th>
			<th>Agencia</th>
		</tr>
		{{range .}}
		<tr>
			<td>{{.ID}}</td>
			<td>{{.Date}}</td>
			<td>{{.Agency}}</td>
		</tr>
		{{end}}
	</table>
</body>
</html>' > <path>/locale/template.html
```

## Uso
Abre tu navegador y dirígete a http://localhost:8000 (o el puerto que hayas especificado en la variable de entorno port).

## Contribución

1. Haz un fork del repositorio.
2. Crea una rama para tu nueva función o corrección de error (git checkout -b feature/nueva-funcion).
3. Haz tus cambios y haz commit (git commit -am 'Agregada nueva función').
4. Haz push a tu rama (git push origin feature/nueva-funcion).
5. Crea un pull request en GitHub y describe tus cambios.

## Licencia

Este proyecto está bajo la Licencia MIT. Consulta el archivo LICENSE para más detalles.

## Créditos

1. Desarrollado por Edgar Javier ↗.
2. Basado en el tutorial Monitor your Postfix email server logs with Golang ↗.

## Contacto

Para cualquier pregunta o comentario, puedes contactar al desarrollador en akosej9208@gmail.com