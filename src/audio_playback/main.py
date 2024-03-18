""" MQTT connected Audio playback program for Blueberry, making use of MPV

Wishes to be provided with {"id", id: str, "audio": audio: str}, where audio is a WAV file, encoded as b64 bytes then decoded into a string, over MQTT to "bloob/{arguments.device_id}/audio_playback/play"

Will respond with {"id": received_id: str, "audio": received_audio: str}. The audio is an exact copy of what was sent. To "bloob/{arguments.device_id}/audio_playback/finished"
"""
import argparse
import subprocess
import asyncio
import aiomqtt
import sys
import re
import json
import base64
import pathlib
import os
import mpv

default_data_path = pathlib.Path(os.environ['HOME']).joinpath(".config/bloob") 
default_temp_path = pathlib.Path("/dev/shm/bloob")

audio_playback_temp_path = default_temp_path.joinpath("audio_playback")

last_audio_file_path = f"{audio_playback_temp_path}/last_played_audio.wav"
if not os.path.exists(audio_playback_temp_path):
	os.makedirs(audio_playback_temp_path)

audio_playback_system = mpv.MPV()

arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('--host', default="localhost")
arg_parser.add_argument('--port', default=1883)
arg_parser.add_argument('--user')
arg_parser.add_argument('--pass')
arg_parser.add_argument('--device-id', default="test")
arguments = arg_parser.parse_args()

def play(audio):
	## Save last played audio to tmp for debugging
	with open(last_audio_file_path,'wb+') as audio_file:
		#Encoding is like this because the string must first be encoded back into the base64 bytes format, then decoded again, this time as b64, into the original bytes.
		audio_file.write(base64.b64decode(audio.encode()))

	print("Playing received audio")
	audio_playback_system.play(last_audio_file_path)

async def connect():
	async with aiomqtt.Client(arguments.host) as client:
		# await client.subscribe(f"bloob/{arguments.device_id}/audio_recorder/finished") # This is for testing, it'll automatically play what the TTS says
		await client.subscribe(f"bloob/{arguments.device_id}/audio_playback/play")
		async for message in client.messages:
			try:
				message_payload = json.loads(message.payload.decode())
				if(message_payload.get('audio') != None and message_payload.get('id') != None):
					play(message_payload["audio"])

					await client.publish(f"bloob/{arguments.device_id}/audio_playback/finished", json.dumps({"id": message_payload.get('id'), "audio":message_payload["audio"]}))
			except:
				print("Error with payload.")

if __name__ == "__main__":
	asyncio.run(connect())