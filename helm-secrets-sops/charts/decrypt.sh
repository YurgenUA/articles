ls -1 secrets/ | while read file;
  do sops -d --pgp $GPG_KEY_ID secrets/$file > secrets-decrypted/${file:0:(-5)}.decrypted.yaml

done