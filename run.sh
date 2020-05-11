echo "ATTENTION: Database is en Europe/Sofia and images are considering America/Sao_Paulo\n\n"

echo "Considering: both loc and db variable" 

curl -XPOST localhost:8101/tz 

echo "\n\n Considering: only loc" 

curl -XPOST localhost:8102/tz 

echo "\n\n Considering: only DB variable" 

curl -XPOST localhost:8103/tz 

echo "\n\n Considering: using none configuration" 

curl -XPOST localhost:8104/tz 