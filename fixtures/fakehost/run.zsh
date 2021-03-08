ssh-keygen -f "$HOME/.ssh/known_hosts" -R "[localhost]:2222"
docker stop openssh-server || true
docker rm openssh-server || true
docker run -it \
  --name=openssh-server \
  --hostname=openssh-server `#optional` \
  -e PUID=1000 \
  -e PGID=1000 \
  -e TZ=Europe/London \
  -e PUBLIC_KEY='ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCoFJkfmA8vG58i8RB3kwuSNoMMXjUT8TC276nHC1e4BD/nHhYAT8ddm61XI5vAJ4N+kVy94wUlhe9K3m6VOMCGWDn1zzX4wdWg297Fkq2kV5Dss/ABTj2aedoKOZisv3qkb82DNh2rfJOPv3FscNd5gRRYjboQFuQAF6qyi++u3YR/LuL8yNACEzetAxUtmVw93pOhvPj8B21ZD6iPTmUBhnvn4m3IfNcImJk4z022MoqW6EZdBJO3xuwq92Uaoe64lYjfOsSRteqZdfkrIci0G4DK/RrFqySq0FMSQnRNYjJTs2ysnyIfW4+oBKSFiniaU9KmlMszAKB/MVmiWeHcZBqEVXEeaRfp+Lm1MjfW6ly8UhloF7LTGCKXeFnrf2CI1Hghuhd1hW2fS4f9partx6luCCNincUhUtxpvChDxNkdErEWMzmYl4d2pUP3Up9w8SdIsBxPSuChTY3fLgjO/ms9IcKPSL1DY1BABo9PzfiaOq0FnZLL4PuytTGU+kftEZz5KaHxaJjxzjOzCBKXzckreNndt11jnQWWZr6sL06w0ppHAdA6093pGokpM7Z5+atnTmaND1NC+5WHlv4Oe61a2LJ2rLF2lJ5QeX8ugJ5mcWSOpSxwX4qdLteXO0N/VIrKjM0VId3EqsqWp4y7gkQ12jpCOSVAW2w/Lgo7oQ== parroit@andrea-XPS' \
  -e SUDO_ACCESS=true `#optional` \
  -e USER_NAME=andrea.parodi \
  -p 2222:2222 \
meteocima/fake-ssh-host
