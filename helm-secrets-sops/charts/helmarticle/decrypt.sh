mkdir -p $1/secrets-decrypted
ls -1 $1/secrets/ | while read file;
  do sops -d --pgp $2 $1/secrets/$file > $1/secrets-decrypted/${file:0:(-5)}.decrypted.yaml

done