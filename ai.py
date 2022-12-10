import openai
import os
from argparse import ArgumentParser
from dotenv import load_dotenv

load_dotenv()

openai.api_key = os.getenv('OpenApiKey')

start_sequence = "\nAI:"
restart_sequence = "\nHuman: "


def read_history(history):
    with open(history, 'r') as h:
        return h.read()


def write_history(history, text):
    with open(history, 'w') as h:
        h.write(text)


if __name__ == '__main__':
    parser = ArgumentParser()
    parser.add_argument('-q', '--question', type=str, help="what question do you want to ask", default=None, required=True)
    args = parser.parse_args()

    p = read_history("./history.txt")
    # user_input = input("Human: ")
    user_input = args.question
    p += user_input
    response = openai.Completion.create(
        model="text-davinci-003",
        prompt=p,
        temperature=0.9,
        max_tokens=2000,
        top_p=1,
        frequency_penalty=0,
        presence_penalty=0.6,
        stop=[" Human:", " AI:"]
    )

    AI_resp = response['choices'][0]['text']
    print(AI_resp.strip())
    p += f"{AI_resp}\nHuman: "
    
    is_write_history = os.getenv('WriteHistory')
    if is_write_history:
        write_history('./history.txt', p)
