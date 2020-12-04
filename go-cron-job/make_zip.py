import zipfile
import time

#https://stackoverflow.com/questions/57964626/permissions-denied-when-trying-to-invoke-go-aws-lambda-function

def make_info(filename):
    info = zipfile.ZipInfo(filename)
    info.date_time = time.localtime()
    info.external_attr = 0x81ed0000
    info.create_system = 3
    return info

zip_source = zipfile.ZipFile("bin.zip")
zip_file = zipfile.ZipFile("dest_file.zip", "w", zipfile.ZIP_DEFLATED)

for cur in zip_source.infolist():
    zip_file.writestr(make_info(cur.filename), zip_source.open(cur.filename).read(), zipfile.ZIP_DEFLATED)

zip_file.close()