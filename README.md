# GO-SOCKS
Socks proxy in Golang

# Senario
![image](https://github.com/user-attachments/assets/3a49392d-6a23-444f-87fb-97dd03d88823)

Let's say you want to access Server C. However, you could only access it via Server B. To achieve this, we will run a SOCKS proxy on Server B and use that to pivot internally.

# Compile
```
go build -ldflags="-s -w" -o go-socks.exe main.go
```
# Update Proxychains Configuration
![image](https://github.com/user-attachments/assets/211b0d79-5738-4edd-9488-4b341dacb148)
```
socks5 Server_B_ip 1080
```
Now run go-socks.exe on Server B; it will open a SOCKS proxy on port 1080 on Server B.

# Disclaimer
The author is not responsible for unauthorized use of this tool. Use responsibly and ensure compliance with legal and ethical standards.
